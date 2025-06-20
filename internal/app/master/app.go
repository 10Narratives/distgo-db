package masterapp

import (
	"log/slog"
	"strconv"

	mastergrpc "github.com/10Narratives/distgo-db/internal/app/master/grpc"
	mastercfg "github.com/10Narratives/distgo-db/internal/config/master"
	databaseapi "github.com/10Narratives/distgo-db/internal/grpc/master/data/database"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type App struct {
	GRPC *mastergrpc.App
}

func New(log *slog.Logger, cfg mastercfg.Config) *App {
	grpcApp := mastergrpc.New(log, cfg.GRPC.Port)

	target := strconv.Itoa(cfg.Worker.Port)
	conn, err := grpc.NewClient("localhost:"+target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("failed to connect to database service", "error", err)
		panic(err)
	}

	databaseClient := dbv1.NewDatabaseServiceClient(conn)
	databaseapi.Register(grpcApp.GRPCServer, databaseClient)

	return &App{GRPC: grpcApp}
}
