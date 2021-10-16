package main

import (
	"fail2connect/config"
	"fail2connect/watcher"
	"fmt"
	"github.com/akamensky/argparse"
	"log"
	"os"
	"time"
)

func main() {
	/// ---------
	/// Parse commandline arguments
	/// ---------
	parser := argparse.NewParser("fail2connect", "Ban connections that fail to connect")
	configFileLocation := parser.String("c", "config", &argparse.Options{Help: "Location of the config file", Default: "config.json"})

	err := parser.Parse(os.Args)
	if err != nil {
		fmt.Print(parser.Usage(err))
		os.Exit(0)
	}

	/// ---------
	/// Parse config file
	/// ---------
	configuration, err := config.ParseConfig(*configFileLocation)

	if err != nil {
		log.Fatalln(err)
	}

	/// ---------
	/// Setup watchers
	/// ---------
	var watchers []*watcher.Watcher

	for _, watcherConfig := range configuration.Watchers {
		if watcherConfig.Enabled {
			watchers = append(watchers, watcher.New(watcherConfig))
		}
	}

	/// ---------
	/// Review connections
	/// ---------
	for true {
		for _, watch := range watchers {
			watch.Review()
		}

		time.Sleep(1 * time.Second)
	}
}
