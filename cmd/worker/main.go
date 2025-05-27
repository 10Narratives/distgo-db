package main

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"

	workerapp "github.com/10Narratives/distgo-db/internal/app/worker"
	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
	"github.com/10Narratives/distgo-db/internal/lib/logger/sl"
)

func main() {
	cfg := workercfg.MustLoad()
	fmt.Println(cfg)

	log := sl.MustLogger(
		sl.WithLevel(cfg.Logging.Level),
		sl.WithFormat(cfg.Logging.Format),
		sl.WithOutput(cfg.Logging.Output),
		sl.WithFileOptions(cfg.Logging.FilePath, cfg.Logging.MaxSize, cfg.Logging.MaxAge, cfg.Logging.Compress),
	)

	log.Info("Worker Node is online")

	application := workerapp.New(log, cfg)

	go func() {
		application.GRPCServer.MustRun()
	}()

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCServer.Stop()
	log.Info("Worker Node stopped")
}
