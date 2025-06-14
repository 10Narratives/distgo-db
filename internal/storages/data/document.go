package datastorage

import (
	"context"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
)

func (s *Storage) CreateDocument(ctx context.Context, name string, value string) (documentmodels.Document, error) {
	panic("unimplemented")
}

func (s *Storage) DeleteDocument(ctx context.Context, name string) error {
	panic("unimplemented")
}

func (s *Storage) Document(ctx context.Context, name string) (documentmodels.Document, error) {
	panic("unimplemented")
}

func (s *Storage) Documents(ctx context.Context, parent string) []documentmodels.Document {
	panic("unimplemented")
}

func (s *Storage) UpdateDocument(ctx context.Context, document documentmodels.Document) error {
	panic("unimplemented")
}
