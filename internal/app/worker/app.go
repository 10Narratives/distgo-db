package workerapp

import (
	"log/slog"

	workergrpc "github.com/10Narratives/distgo-db/internal/app/worker/grpc"
	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
	documentstore "github.com/10Narratives/distgo-db/internal/storages/worker/document"
	walstore "github.com/10Narratives/distgo-db/internal/storages/worker/wal"

	documentsrv "github.com/10Narratives/distgo-db/internal/services/worker/document"
)

type App struct {
	GRPC *workergrpc.App
}

func New(log *slog.Logger, cfg workercfg.Config) *App {
	documentStorage := documentstore.NewStorage()
	walStorage, err := walstore.New("logs/" + cfg.Name + ".log")
	if err != nil {
		panic(err.Error())
	}

	documentService := documentsrv.New(documentStorage, walStorage)

	grpcApp := workergrpc.New(log, documentService, cfg.GRPC.Port)
	return &App{GRPC: grpcApp}
}
