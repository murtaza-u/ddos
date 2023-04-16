package main

import (
	"log"
	"os"

	"github.com/murtaza-u/ddos/apisrv"
	"github.com/murtaza-u/ddos/conf"
)

const defaultConfig = "/etc/ddos/config.yaml"

func main() {
	path := os.Getenv("DDOS_APISRV_CONFIG")
	if path == "" {
		path = defaultConfig
	}

	c, err := conf.New(path)
	if err != nil {
		log.Fatal(err)
	}

	err = c.Validate()
	if err != nil {
		log.Fatalf("failed to validate config %s", err.Error())
	}

	err = apisrv.Start(c)
	if err != nil {
		log.Fatal(err)
	}
}
