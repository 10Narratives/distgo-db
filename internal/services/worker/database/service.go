package databasesrv

import (
	"context"
	"log/slog"

	"github.com/10Narratives/distgo-db/internal/models"
)

type DatabaseService struct {
	log *slog.Logger
}

func New(log *slog.Logger) *DatabaseService {
	return &DatabaseService{log: log}
}

func (srv *DatabaseService) ListDocuments(ctx context.Context, collection string) ([]models.Document, error) {
	panic("ii")
}

func (srv *DatabaseService) GetDocument(ctx context.Context, collection string, documentID string) (models.Document, error) {
	panic("ii")
}

func (srv *DatabaseService) CreateDocument(ctx context.Context, collection string, data map[string]any) (models.Document, error) {
	panic("ii")
}

func (srv *DatabaseService) UpdateDocument(ctx context.Context, collection string, update map[string]any) (models.Document, error) {
	panic("ii")
}

func (srv *DatabaseService) DeleteDocument(ctx context.Context, collection, documentID string) (bool, error) {
	panic("ii")
}
