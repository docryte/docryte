package main

import (
	"log"

	"docryte/src/api"
	"docryte/src/config"
)

func main() {
	cfg, err := config.Get()
	if err != nil {
		log.Fatal(err.Error())
	}
	api.Init(cfg)
}
