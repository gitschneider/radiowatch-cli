package main

import (
	"github.com/schnaidar/stationcrawler"
	"github.com/schnaidar/radiowatch"
	"github.com/codegangsta/cli"
	"fmt"
	"os"
	"io/ioutil"
	"encoding/json"
	log "github.com/Sirupsen/logrus"
)

type config struct {
	Database string `json:"database"`
	Host     string `json:"host"`
	Password string `json:"password"`
	Port     string `json:"port"`
	User     string `json:"user"`
}

func main() {
	log.SetFormatter(&log.TextFormatter{ForceColors:true})
	app := cli.NewApp()

	app.Action = func(c *cli.Context) {
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
			stationcrawler.NewAntenne(),
			stationcrawler.NewMdrJump(),
		})

		watcher.StartCrawling()
		channel := make(chan bool)
		<-channel
	}

	app.Run(os.Args)
}
