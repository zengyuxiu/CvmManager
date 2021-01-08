package lxc_api

import (
	"database/sql"
	"github.com/lxc/lxd/client"
	"github.com/lxc/lxd/shared/api"
	_ "github.com/mattn/go-sqlite3"
	log "github.com/sirupsen/logrus"
	"time"
)

func LxcStatus(sqliteDatabase *sql.DB) error {
	c, err := lxd.ConnectLXDUnix("",nil)
	if err != nil{
		log.Error(err)
	}
	containers,err := c.GetContainers()
	for _,container := range containers{
		state,_,err := c.GetInstanceState(container.Name)
		if err != nil{
			log.Error(err)
		}else{
			InsertStatLxc(sqliteDatabase,state,container.Name,time.Now().UTC())
		}
	}
	return nil
}

//goland:noinspection SqlResolve
func InsertStatLxc(db *sql.DB,stat *api.InstanceState,name string,time time.Time)  {
	var InsertSQL = `INSERT INTO lxc_stats(
                    name,
                    cpu_usage,
                    memory_usage,
                    memory_max_usage,
                    time_stamp
                    ) VALUES (?, ?, ?, ?, ?)`
	statement, err := db.Prepare(InsertSQL)
	if err != nil {
		log.Error(err)
	}
	_, err = statement.Exec(name,stat.CPU.Usage,
									stat.Memory.Usage,stat.Memory.UsagePeak,
									time,
	)
	if err != nil {
		log.Error(err)
	}
}
