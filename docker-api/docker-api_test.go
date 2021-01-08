package docker_api

import (
	"database/sql"
	"testing"
	log "github.com/sirupsen/logrus"
)

func TestStatus(t *testing.T)  {
	var sqliteDatabase, err = sql.Open("sqlite3", "../db-api/CvmStats.db")
	if err != nil{
		log.Fatal(err)
	}
	defer sqliteDatabase.Close()
	err = DockerStatus(sqliteDatabase)
	print(err)
}

