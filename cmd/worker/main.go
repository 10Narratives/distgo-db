package main

import (
	"os"
	"os/signal"
	"syscall"

	workerapp "github.com/10Narratives/distgo-db/internal/app/worker"
	"github.com/10Narratives/distgo-db/internal/config"
	"github.com/10Narratives/distgo-db/internal/lib/logger/sl"
)

func main() {
	cfg := config.MustLoadWorker()

	log := sl.MustLogger(
		sl.WithLevel(cfg.Logging.Level),
		sl.WithFormat(cfg.Logging.Format),
		sl.WithOutput(cfg.Logging.Output),
		sl.WithFileOptions(cfg.Logging.FilePath, cfg.Logging.MaxSize, cfg.Logging.MaxAge, cfg.Logging.Compress),
	)

	log.Info("Worker node started")

	application := workerapp.New(log, cfg)

	go func() {
		application.GRPCSrv.MustRun()
	}()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	application.GRPCSrv.Stop()

	log.Info("Worker node stopped")
}
