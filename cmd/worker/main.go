package main

import (
	"fmt"

	"githib.com/10Narratives/distgo-db/internal/worker/config"
	"githib.com/10Narratives/distgo-db/internal/worker/logging"
)

func main() {
	cfg := config.MustLoad()
	fmt.Println(cfg)

	log := logging.NewLogger(cfg.Logging.Format, cfg.Logging.Level)
	log.Info("starting worker node #1")

	// TODO: initialize storage app

	// TODO: initialize gRPC-server app

	// TODO: initialize worker app

	log.Info("worker node #1 stopped")
}
