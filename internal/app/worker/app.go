package workerapp

import (
	"log/slog"

	workergrpc "github.com/10Narratives/distgo-db/internal/app/worker/grpc"
	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
	collectionsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/collection"
	databasesrv "github.com/10Narratives/distgo-db/internal/services/worker/data/database"
	documentsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/document"
	datastorage "github.com/10Narratives/distgo-db/internal/storages/data"
)

type App struct {
	GRPC *workergrpc.App
}

func New(log *slog.Logger, cfg workercfg.Config) *App {
	storage := datastorage.New()

	databaseSrv := databasesrv.New(storage, nil)
	collectionSrv := collectionsrv.New(storage, nil)
	documentSrv := documentsrv.New(storage, nil)

	grpcApp := workergrpc.New(
		log,
		databaseSrv,
		collectionSrv,
		documentSrv,
		cfg.GRPC.Port,
	)
	return &App{GRPC: grpcApp}
}
