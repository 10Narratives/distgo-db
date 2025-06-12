package documentstore

// import (
// 	"context"

// 	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
// 	documentsrv "github.com/10Narratives/distgo-db/internal/services/worker/document"
// 	"github.com/google/uuid"
// )

// type Storage struct {
// 	collections map[string]*documentmodels.Collection
// }

// var _ documentsrv.DocumentStorage = Storage{}

// func New() *Storage {
// 	return nil
// }

// func (s Storage) Delete(ctx context.Context, collection string, documentID uuid.UUID) error {
// 	panic("unimplemented")
// }

// func (s Storage) Get(ctx context.Context, collection string, documentID uuid.UUID) (documentmodels.Document, error) {
// 	panic("unimplemented")
// }

// func (s Storage) List(ctx context.Context, collection string) ([]documentmodels.Document, error) {
// 	panic("unimplemented")
// }

// func (s Storage) Set(ctx context.Context, collection string, documentID uuid.UUID, content map[string]any) {
// 	panic("unimplemented")
// }
