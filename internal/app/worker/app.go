package workerapp

import (
	"log/slog"

	workergrpc "github.com/10Narratives/distgo-db/internal/app/worker/grpc"
	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
	documentstore "github.com/10Narratives/distgo-db/internal/storages/worker/document"

	documentsrv "github.com/10Narratives/distgo-db/internal/services/worker/document"
)

type App struct {
	GRPC *workergrpc.App
}

func New(log *slog.Logger, cfg workercfg.Config) *App {
	documentStorage := documentstore.New()
	documentService := documentsrv.New(documentStorage)
	grpcApp := workergrpc.New(log, documentService, cfg.GRPC.Port)
	return &App{GRPC: grpcApp}
}
