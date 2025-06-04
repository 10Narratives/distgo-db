package workerapp

import (
	"log/slog"

	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
	databasesrv "github.com/10Narratives/distgo-db/internal/services/worker/database"
	databasestore "github.com/10Narratives/distgo-db/internal/storage/database"
)

type App struct {
	GRPCServer *GRPCApp
}

func New(log *slog.Logger, cfg *workercfg.Config) *App {
	databaseStore := databasestore.New()

	databaseService := databasesrv.New(log, databaseStore)

	grpcApp := NewGRPCApp(log, databaseService, cfg.GRPC.Port)
	return &App{GRPCServer: grpcApp}
}
