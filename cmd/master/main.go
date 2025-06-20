package main

import (
	"os"
	"os/signal"
	"syscall"

	masterapp "github.com/10Narratives/distgo-db/internal/app/master"
	mastercfg "github.com/10Narratives/distgo-db/internal/config/master"
	"github.com/10Narratives/distgo-db/internal/lib/logger/sl"
)

func main() {
	cfg := mastercfg.MustLoad()

	log := sl.MustLogger(
		sl.WithFormat(cfg.Logging.Format),
		sl.WithOutput(cfg.Logging.Output),
		sl.WithLevel(cfg.Logging.Level),
	)
	log.Info("Master is online")

	app := masterapp.New(log, *cfg)

	go func() {
		app.GRPC.MustRun()
	}()

	stop := make(chan os.Signal, 2)
	signal.Notify(stop, syscall.SIGTERM, syscall.SIGINT)

	<-stop

	app.ForgetAllWorkers()
	app.GRPC.Stop()

	log.Info("Master is stopped")
}
