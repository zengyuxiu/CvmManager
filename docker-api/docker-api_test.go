package docker_api

import (
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestStatus(t *testing.T) {
	var sqliteDatabase, err = sql.Open("sqlite3", "../db-api/CvmStats.db")
	if err != nil {
		log.Fatal(err)
	}
	defer sqliteDatabase.Close()
	err = DockerStatus(sqliteDatabase)
	print(err)
}

func TestDockerCreate(t *testing.T) {
	var (
		InstanceNum = 3
		Image       = "debian"
	)
	err := DockerCreate(InstanceNum, Image)
	if err != nil {
		fmt.Printf("OK")
	}
}

func TestDockerDelete(t *testing.T) {
	err := DockerDelete()
	if err != nil {
		fmt.Printf("OK")
	}
}
