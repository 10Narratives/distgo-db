package workerapp

import (
	"log/slog"

	"github.com/10Narratives/distgo-db/internal/config"
)

type App struct {
	GRPCSrv *WorkerGRPCApp
}

func New(log *slog.Logger, cfg *config.WorkerConfig) *App {

	grpcApp := NewGRPCApp(log, nil, cfg.GRPC.Port)

	return &App{GRPCSrv: grpcApp}
}
