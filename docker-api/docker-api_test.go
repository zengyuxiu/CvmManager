package docker_api

import (
	"CvmManager/config"
	"database/sql"
	"fmt"
	log "github.com/sirupsen/logrus"
	"testing"
)

func TestStatus(t *testing.T) {
	var sqliteDatabase, err = sql.Open("sqlite3", "../db-api/testDB.db")
	if err != nil {
		log.Fatal(err)
	}
	defer sqliteDatabase.Close()
	err = DockerStatus(sqliteDatabase)
	print(err)
}

func TestDockerCreate(t *testing.T) {
	var (
		InstanceNum = 100
		Image       = "6d6b00c22231"
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

func TestNetworkCreate(t *testing.T) {
	networkConfig := config.GetConfig()
	networkConfig.OspfdConfig()
	NetworkCreate(networkConfig)
}
