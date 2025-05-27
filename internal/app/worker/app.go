package workerapp

import (
	"log/slog"

	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
)

type App struct {
	GRPCServer *GRPCApp
}

func New(log *slog.Logger, cfg *workercfg.Config) *App {
	grpcApp := NewGRPCApp(log, nil, cfg.GRPC.Port)
	return &App{GRPCServer: grpcApp}
}
