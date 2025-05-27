package main

import (
	"fmt"

	mastercfg "github.com/10Narratives/distgo-db/internal/config/master"
	"github.com/10Narratives/distgo-db/internal/lib/logger/sl"
)

func main() {
	cfg := mastercfg.MustLoad()
	fmt.Println(cfg)

	log := sl.MustLogger(
		sl.WithLevel(cfg.Logging.Level),
		sl.WithFormat(cfg.Logging.Format),
		sl.WithOutput(cfg.Logging.Output),
		sl.WithFileOptions(cfg.Logging.FilePath, cfg.Logging.MaxSize, cfg.Logging.MaxAge, cfg.Logging.Compress),
	)

	log.Info("Master Node is online")

	log.Info("Master Node stopped")
}
