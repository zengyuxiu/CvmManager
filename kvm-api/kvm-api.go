package kvm_api

import (
	"database/sql"
	"fmt"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"libvirt.org/libvirt-go"
	"time"
)

func KvmStatus(sqliteDatabase *sql.DB) error {
	conn, err := libvirt.NewConnect("qemu:///system")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	doms, err := conn.ListAllDomains(libvirt.CONNECT_LIST_DOMAINS_ACTIVE)
	if err != nil {
		log.Fatal(err)
	}
	for _, dom := range doms {
		name, err := dom.GetName()
		info, err := dom.GetInfo()
		if err == nil {
			InsertStatKvm(sqliteDatabase, info, name, time.Now().UTC())
		}
		dom.Free()
	}
	return nil
}

//goland:noinspection SqlResolve
func InsertStatKvm(db *sql.DB, stat *libvirt.DomainInfo, name string, now time.Time) {
	var InsertSQL = `INSERT INTO kvm_stats(
                    Name,
                    memory_usage,
                    memory_max_usage,
                    cpu_usage,
                    cpu_online,
                    time_stamp
                    ) VALUES (?, ?, ?, ?, ?, ?)`
	statement, err := db.Prepare(InsertSQL)
	if err != nil {
		log.Error(err)
	}
	result, err := statement.Exec(name, stat.Memory,
		stat.MaxMem, stat.CpuTime, stat.NrVirtCpu, now,
	)
	if err != nil {
		log.Error(err)
	}
	fmt.Print(result)
}
