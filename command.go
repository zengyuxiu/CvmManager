package main

import (
	dockerapi "CvmManager/docker-api"
	kvmapi "CvmManager/kvm-api"
	lxcapi "CvmManager/lxc-api"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var createCommand = cli.Command{
	Name: "create",
	Usage: "Create instances",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name: "i",
			Usage: "specify instructure",
		},
	},
}
var statusCommand = cli.Command{
	Name: "stat",
	Usage: "Write the instances status to db",
	Flags: []cli.Flag {
		cli.StringFlag{
			Name: "i",
			Usage: "specify instructure",
		},
		cli.BoolFlag{
			Name: "d",
			Usage: "daemon run",
		},
	},
	Action: func(ctx *cli.Context) error {
		var sqliteDatabase, err = sql.Open("sqlite3", "db-api/CvmStats.db")
		if err != nil{
			log.Fatal(err)
		}
		defer sqliteDatabase.Close()
		instructure := ctx.String("i")
		switch instructure {
		case "docker":
			err := dockerapi.DockerStatus(sqliteDatabase)
			if err != nil{
				log.Fatal(err)
			}
			break
		case "lxc":
			err := lxcapi.LxcStatus(sqliteDatabase)
			if err != nil{
				log.Fatal(err)
			}
			break
		case "kvm":
			err := kvmapi.KvmStatus(sqliteDatabase)
			if err != nil{
				log.Fatal(err)
			}
			break
		default:
			log.Error("Unknown Instructure")
			return nil
		}
		return nil
	},
}
var deleteCommand = cli.Command{

}
