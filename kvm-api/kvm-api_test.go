package kvm_api

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestKvmStatus(t *testing.T) {
	var sqliteDatabase, err = sql.Open("sqlite3", "../db-api/CvmStats.db")
	if err != nil {
		log.Fatal(err)
	}
	defer sqliteDatabase.Close()
	if err := KvmStatus(sqliteDatabase); err != nil {
		t.Errorf("Error")
	} else {
		fmt.Print("Hello")
	}
}
