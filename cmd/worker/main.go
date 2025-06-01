package main

import (
	"fmt"

	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
	"github.com/10Narratives/distgo-db/internal/lib/logging/sl"
)

func main() {
	cfg := workercfg.MustLoad()
	fmt.Println(cfg)

	log := sl.MustLogger(
		sl.WithLevel(cfg.Logging.Level),
		sl.WithFormat(cfg.Logging.Format),
		sl.WithOutput(cfg.Logging.Output),
	)
	log.Info("Worker " + cfg.Name + " is online")

	// TODO: initialize gRPC application
	// TODO: initialize application
}
