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
var connectionLimit = flag.Int("l", 0, "Connection limit (0 = unlimited).")

func main() {
	flag.Parse()

	// config
	cfg, err := LoadConfig(*cfgFileName)
	if err != nil {
		log.Fatal(err)
	}
	if cfg.isCreated {
		log.Printf("Created default config file: %s", cfg.configFileName)
	}
	log.Printf("Using config file: %s", cfg.configFileName)

	srv := NewServer(cfg.Server.Address)
	srv.ConnLimit = *connectionLimit

	log.Fatal(srv.Http.ListenAndServe())
}
