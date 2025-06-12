package workerapp

import (
	"log/slog"

	workergrpc "github.com/10Narratives/distgo-db/internal/app/worker/grpc"
	workercfg "github.com/10Narratives/distgo-db/internal/config/worker"
)

type App struct {
	GRPC *workergrpc.App
}

func New(log *slog.Logger, cfg workercfg.Config) *App {
	// documentStorage := documentstore.New()
	// walStorage, err := walstore.New("logs/" + cfg.Name + ".log")
	// if err != nil {
	// 	panic(err.Error())
	// }

	// documentService := documentsrv.New(documentStorage, walStorage)

	grpcApp := workergrpc.New(log, cfg.GRPC.Port)
	return &App{GRPC: grpcApp}
}
