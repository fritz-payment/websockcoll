package main

import (
	"flag"
	"log"
)

const (
	AppName    = "coll_server"
	AppVersion = "0.1"
)

// Cmd line flags
var cfgFileName = flag.String("c", "", "Config file name to read.")

func main() {
	flag.Parse()
	cfg, err := LoadConfig(*cfgFileName)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.isCreated {
		log.Printf("Created default config file: %s", cfg.configFileName)
	}
	log.Printf("Using config file: %s", cfg.configFileName)
}
