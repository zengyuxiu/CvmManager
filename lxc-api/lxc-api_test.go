package lxc_api

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestLxcStatus(t *testing.T) {
	var sqliteDatabase, err = sql.Open("sqlite3", "../db-api/testDB.db")
	if err != nil {
		log.Fatal(err)
	}
	defer sqliteDatabase.Close()
	if err := LxcStatus(sqliteDatabase); err != nil {
		t.Errorf("Error")
	} else {
		fmt.Print("Hello")
	}
}

func TestLxcCreate(t *testing.T) {
	var (
		InstanceNum = 5
		Fingerprint = "9a55859f4761"
	)
	if err := LxcCreate(InstanceNum, Fingerprint); err != nil {
		t.Errorf("Error")
		log.Fatal(err)
	}
}

func TestLxcDelete(t *testing.T) {
	if err := LxcDelete(); err != nil {
		t.Errorf("Error")
		log.Fatal(err)
	}
}

func TestGetRandomString(t *testing.T) {
	text := GetRandomString(8)
	fmt.Printf(string(text))
}
