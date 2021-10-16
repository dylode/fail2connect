package main

import (
	"fmt"
	"github.com/akamensky/argparse"
	"log"
	"os"
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
	config, err := ParseConfig(*configFileLocation)

	if err != nil {
		log.Fatalln(err)
	}

	/// ---------
	/// Setup watcher
	/// ---------
	watchers := config.Watchers
	fmt.Print(watchers[0].BanCommand)


}
