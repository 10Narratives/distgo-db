package workerapp

import (
	"log/slog"

	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
	databasesrv "github.com/10Narratives/distgo-db/internal/services/worker/database"
)

type App struct {
	GRPCServer *GRPCApp
}

func New(log *slog.Logger, cfg *workercfg.Config) *App {

	databaseSrv := databasesrv.New(log)

	grpcApp := NewGRPCApp(log, databaseSrv, cfg.GRPC.Port)
	return &App{GRPCServer: grpcApp}
}
