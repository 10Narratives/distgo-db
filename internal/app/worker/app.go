package workerapp

import (
	"context"
	"log/slog"

	workergrpc "github.com/10Narratives/distgo-db/internal/app/worker/grpc"
	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
	clustergrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/cluster"
	collectiongrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/collection"
	databasegrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/database"
	documentgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/document"
	transactiongrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/transaction"
	walgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/wal"
	clustersrv "github.com/10Narratives/distgo-db/internal/services/worker/cluster"
	collectionsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/collection"
	databasesrv "github.com/10Narratives/distgo-db/internal/services/worker/data/database"
	documentsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/document"
	transactionsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/transaction"
	walsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/wal"
	datastorage "github.com/10Narratives/distgo-db/internal/storages/data"
	transactionstorage "github.com/10Narratives/distgo-db/internal/storages/transaction"
	walstorage "github.com/10Narratives/distgo-db/internal/storages/wal"
)

type App struct {
	GRPC           *workergrpc.App
	ClusterService *clustersrv.Service
}

func New(log *slog.Logger, cfg workercfg.Config) *App {
	grpcApp := workergrpc.New(log, cfg.GRPC.Port, cfg.Master.Port, cfg.Name)

	walStorage, err := walstorage.New(cfg.WAL.Path)
	if err != nil {
		panic("cannot initialize wal storage")
	}

	dataStorage := datastorage.New()
	err = dataStorage.RecoverFromFile(context.Background(), cfg.WAL.Path)
	if err != nil {
		panic("cannot recover data storage from wal log")
	}

	txStorage := transactionstorage.New()

	walService := walsrv.New(walStorage)
	walgrpc.Register(grpcApp.GRPCServer, walService)

	databaseSrv := databasesrv.New(dataStorage, walService)
	databasegrpc.Register(grpcApp.GRPCServer, databaseSrv)

	collectionSrv := collectionsrv.New(dataStorage, walService)
	collectiongrpc.Register(grpcApp.GRPCServer, collectionSrv)

	documentSrv := documentsrv.New(dataStorage, walService)
	documentgrpc.Register(grpcApp.GRPCServer, documentSrv)

	txService := transactionsrv.New(txStorage, dataStorage)
	transactiongrpc.Register(grpcApp.GRPCServer, txService)

	clusterService := clustersrv.New(cfg.Master.Port)
	clustergrpc.Register(grpcApp.GRPCServer, clusterService)

	return &App{GRPC: grpcApp, ClusterService: clusterService}
}
