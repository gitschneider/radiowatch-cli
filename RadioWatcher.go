package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	log "github.com/Sirupsen/logrus"
	"github.com/codegangsta/cli"
	"github.com/gitschneider/radiowatch"
	"github.com/gitschneider/stationcrawler"
)

type config struct {
	Database string `json:"database"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     string `json:"port"`
	User     string `json:"user"`
}

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors: true})
	log.SetLevel(log.ErrorLevel)
	app := cli.NewApp()

	app.Flags = []cli.Flag{
		cli.BoolFlag{
			Name:  "verbose",
			Usage: "Outputs debug information",
		},
	}

	app.Action = func(c *cli.Context) {
		if c.Bool("verbose") {
			log.SetLevel(log.DebugLevel)
		}
		var cfg config
		file, err := ioutil.ReadFile("radiowatch.json")
		if err != nil {
			log.WithField(
				"message",
				err.Error(),
			).Fatal("Error when reading config")
		}
		err = json.Unmarshal(file, &cfg)
		if err != nil {
			fmt.Fprintln(os.Stderr, err.Error())
			os.Exit(2)
			log.WithField(
				"message",
				err.Error(),
			).Fatal("Error when Unmarshal config")
		}

		writer := radiowatch.NewMysqlWriter(cfg.User, cfg.Password, cfg.Host, cfg.Port, cfg.Database)
		watcher := radiowatch.NewWatcher(writer)
		watcher.SetInterval("20s")

		watcher.AddCrawlers([]radiowatch.Crawler{
			stationcrawler.NewNjoy(),
			stationcrawler.NewNdr2(),
			stationcrawler.NewDasDing(),
			stationcrawler.NewHr3(),
			stationcrawler.NewYouFm(),
			stationcrawler.NewFfn(),
			stationcrawler.NewMdrJump(),
		})

		watcher.StartCrawling()
		channel := make(chan bool)
		<-channel
	}

	app.Run(os.Args)
}
