package lxc_api

import (
	"database/sql"
	"fmt"
	"testing"
	log "github.com/sirupsen/logrus"
)

func TestLxcStatus(t *testing.T)  {
	var sqliteDatabase, err = sql.Open("sqlite3", "../db-api/CvmStats.db")
	if err != nil{
		log.Fatal(err)
	}
	defer sqliteDatabase.Close()
	if err := LxcStatus(sqliteDatabase) ; err != nil{
		t.Errorf("Error")
	}else {
		fmt.Print("Hello")
	}
}
