package main

import (
	log "github.com/sirupsen/logrus"
	"github.com/urfave/cli"
	"os"
)

const Usage = "A CLI Container & VM manager"

func main() {
	app := cli.NewApp()
	app.Name = "ldkmngr"
	app.Usage = Usage
	app.Commands = []cli.Command{
		createCommand,
		statusCommand,
		deleteCommand,
	}
	app.Before = func(context *cli.Context) error {
		// Log as JSON instead of the default ASCII formatter.
		log.SetFormatter(&log.JSONFormatter{})

		log.SetOutput(os.Stdout)
		return nil
	}

	if err := app.Run(os.Args); err != nil {
		log.Fatal(err)
	}
}
