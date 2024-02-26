package main

import (
	"git.spbec-mining.ru/arxon31/sambaMW/internal/app"
	"git.spbec-mining.ru/arxon31/sambaMW/internal/config"
	"log"
)

func main() {
	cfg, err := config.MustLoad()
	if err != nil {
		log.Fatalf("can not read config file: %s", err)
	}

	app.Run(cfg)

}
