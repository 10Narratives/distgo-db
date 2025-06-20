package masterapp

import (
	"log/slog"

	workergrpc "github.com/10Narratives/distgo-db/internal/app/worker/grpc"
	mastercfg "github.com/10Narratives/distgo-db/internal/config/master"
)

type App struct {
	GRPC *workergrpc.App
}

func New(log *slog.Logger, cfg mastercfg.Config) *App {
	grpcApp := workergrpc.New(log, cfg.GRPC.Port)
	return &App{GRPC: grpcApp}
}
