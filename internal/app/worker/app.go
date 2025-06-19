package workerapp

import (
	"context"
	"log/slog"

	workergrpc "github.com/10Narratives/distgo-db/internal/app/worker/grpc"
	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
	collectiongrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/collection"
	databasegrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/database"
	documentgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/document"
	walgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/wal"
	collectionsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/collection"
	databasesrv "github.com/10Narratives/distgo-db/internal/services/worker/data/database"
	documentsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/document"
	walsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/wal"
	datastorage "github.com/10Narratives/distgo-db/internal/storages/data"
	walstorage "github.com/10Narratives/distgo-db/internal/storages/wal"
)

type App struct {
	GRPC *workergrpc.App
}

func New(log *slog.Logger, cfg workercfg.Config) *App {
	grpcApp := workergrpc.New(log, cfg.GRPC.Port)

	walStorage, err := walstorage.New(cfg.WAL.Path)
	if err != nil {
		panic("cannot initialize wal storage")
	}

	dataStorage := datastorage.New()
	err = dataStorage.RecoverFromFile(context.Background(), cfg.WAL.Path)
	if err != nil {
		panic("cannot recover data storage from wal log")
	}

	walService := walsrv.New(walStorage)
	walgrpc.Register(grpcApp.GRPCServer, walService)

	databaseSrv := databasesrv.New(dataStorage, walService)
	databasegrpc.Register(grpcApp.GRPCServer, databaseSrv)

	collectionSrv := collectionsrv.New(dataStorage, walService)
	collectiongrpc.Register(grpcApp.GRPCServer, collectionSrv)

	documentSrv := documentsrv.New(dataStorage, walService)
	documentgrpc.Register(grpcApp.GRPCServer, documentSrv)

	return &App{GRPC: grpcApp}
}
