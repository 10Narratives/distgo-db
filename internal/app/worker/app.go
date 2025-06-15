package workerapp

import (
	"log/slog"

	workergrpc "github.com/10Narratives/distgo-db/internal/app/worker/grpc"
	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
	databasesrv "github.com/10Narratives/distgo-db/internal/services/worker/data/database"
	datastorage "github.com/10Narratives/distgo-db/internal/storages/data"
)

type App struct {
	GRPC *workergrpc.App
}

func New(log *slog.Logger, cfg workercfg.Config) *App {
	storage := datastorage.New()

	databaseSrv := databasesrv.New(storage)

	grpcApp := workergrpc.New(log, databaseSrv, cfg.GRPC.Port)
	return &App{GRPC: grpcApp}
}
