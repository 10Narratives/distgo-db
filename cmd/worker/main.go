package main

import (
	"os"
	"os/signal"
	"syscall"

	workerapp "github.com/10Narratives/distgo-db/internal/app/worker"
	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
	"github.com/10Narratives/distgo-db/internal/lib/logger/sl"
)

func main() {
	cfg := workercfg.MustLoad()

	log := sl.MustLogger(
		sl.WithFormat(cfg.Logging.Format),
		sl.WithOutput(cfg.Logging.Output),
		sl.WithLevel(cfg.Logging.Level),
	)
	log.Info(cfg.Name + " is online")

	app := workerapp.New(log, *cfg)
	app.ClusterService.MustRegister(cfg.GRPC.Port, cfg.Name)

	go func() {
		app.GRPC.MustRun()
	}()

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	err := app.ClusterService.Unregister()
	if err != nil {
		log.Error(err.Error())
	}

	app.GRPC.Stop()
	log.Info(cfg.Name + " is stopped")
}
