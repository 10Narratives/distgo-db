package masterapp

import (
	"log/slog"

	mastergrpc "github.com/10Narratives/distgo-db/internal/app/master/grpc"
	mastercfg "github.com/10Narratives/distgo-db/internal/config/master"
	clusterapi "github.com/10Narratives/distgo-db/internal/grpc/master/cluster"
	clustersrv "github.com/10Narratives/distgo-db/internal/services/master/cluster"
	clusterstorage "github.com/10Narratives/distgo-db/internal/storages/cluster"
)

type App struct {
	GRPC *mastergrpc.App
}

func New(log *slog.Logger, cfg mastercfg.Config) *App {
	grpcApp := mastergrpc.New(log, cfg.GRPC.Port)

	clusterStorage := clusterstorage.New()

	clusterService := clustersrv.New(clusterStorage)
	clusterapi.Register(grpcApp.GRPCServer, clusterService)

	return &App{GRPC: grpcApp}
}
