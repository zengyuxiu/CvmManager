package main

import (
	db_api "CvmManager/db-api"
	dockerapi "CvmManager/docker-api"
	kvmapi "CvmManager/kvm-api"
	lxcapi "CvmManager/lxc-api"
	"database/sql"
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
)

var initCommand = cli.Command{
	Name:  "init",
	Usage: "Create DB and Tables",
	Action: func(ctx *cli.Context) {
		db_api.Init_DB()
		var sqliteDatabase, err = sql.Open("sqlite3", "stats-database.db")
		if err != nil {
			log.Fatal(err)
		}
		defer sqliteDatabase.Close()
		db_api.CreateTable(sqliteDatabase)
	},
}
var createCommand = cli.Command{
	Name:  "create",
	Usage: "Create instances",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "i",
			Usage: "specify instructure \"lxc\" \"docker\"",
		},
		cli.IntFlag{
			Name:  "n",
			Usage: "Instance Number",
		},
		cli.StringFlag{
			Name:  "s",
			Usage: "source image's ID(Docker)/Fingerprint(Lxc)",
		},
	},
	Action: func(ctx *cli.Context) error {
		Instructure := ctx.String("i")
		InstanceNum := ctx.Int("n")
		Image := ctx.String("s")

		switch Instructure {
		case "docker":
			err := dockerapi.DockerCreate(InstanceNum, Image)
			if err != nil {
				log.Fatal(err)
			}
			break
		case "lxc":
			err := lxcapi.LxcCreate(InstanceNum, Image)
			if err != nil {
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
var statusCommand = cli.Command{
	Name:  "stat",
	Usage: "Write the instances status to db",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "i",
			Usage: "specify instructure \"lxc\" \"docker\"",
		},
		/*		TODO
				cli.BoolFlag{
					Name: "d",
					Usage: "daemon run",
				},*/
	},
	Action: func(ctx *cli.Context) error {
		var sqliteDatabase, err = sql.Open("sqlite3", "stats-database.db")
		if err != nil {
			log.Fatal(err)
		}
		defer sqliteDatabase.Close()
		instructure := ctx.String("i")
		switch instructure {
		case "docker":
			err := dockerapi.DockerStatus(sqliteDatabase)
			if err != nil {
				log.Fatal(err)
			}
			break
		case "lxc":
			err := lxcapi.LxcStatus(sqliteDatabase)
			if err != nil {
				log.Fatal(err)
			}
			break
		case "kvm":
			err := kvmapi.KvmStatus(sqliteDatabase)
			if err != nil {
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
	Name:  "stat",
	Usage: "Delete Test Instance",
	Flags: []cli.Flag{
		cli.StringFlag{
			Name:  "i",
			Usage: "specify instructure \"lxd\" \"docker\"",
		},
	},
	Action: func(ctx *cli.Context) error {
		Instructure := ctx.String("i")
		switch Instructure {
		case "docker":
			err := dockerapi.DockerDelete()
			if err != nil {
				log.Fatal(err)
			}
			break
		case "lxc":
			err := lxcapi.LxcDelete()
			if err != nil {
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
