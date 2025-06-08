package documentsrv

import (
	"context"
	"encoding/json"
	"time"

	documentgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/document"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/wal"
	"github.com/google/uuid"
)

//go:generate mockery --name DocumentStorage --output ./mocks/
type DocumentStorage interface {
	Get(ctx context.Context, collection string, documentID uuid.UUID) (documentmodels.Document, error)
	Set(ctx context.Context, collection string, documentID uuid.UUID, content map[string]any)
	List(ctx context.Context, collection string) ([]documentmodels.Document, error)
	Delete(ctx context.Context, collection string, documentID uuid.UUID) error
	Replace(ctx context.Context, collection string, documentID uuid.UUID, content map[string]any) (documentmodels.Document, error)
}

//go:generate mockery --name WALStorage --output ./mocks/
type WALStorage interface {
	Write(entry walmodels.Entry) error
	Replay(handler func(walmodels.Entry) error) error
}

type Service struct {
	documentStorage DocumentStorage
	walStorage      WALStorage
}

func New(
	documentStorage DocumentStorage,
	walStorage WALStorage) *Service {
	service := &Service{
		documentStorage: documentStorage,
		walStorage:      walStorage,
	}

	err := service.walStorage.Replay(func(entry walmodels.Entry) error {
		var record walmodels.Record
		if err := json.Unmarshal(entry, &record); err != nil {
			return err
		}

		switch record.Op {
		case documentmodels.OpCreate:
			service.documentStorage.Set(context.Background(), record.Collection, record.DocumentID, record.Content)
		case documentmodels.OpUpdate:
			_, err := service.documentStorage.Replace(context.Background(), record.Collection, record.DocumentID, record.Content)
			if err != nil {
				return err
			}
		case documentmodels.OpDelete:
			err := service.documentStorage.Delete(context.Background(), record.Collection, record.DocumentID)
			if err != nil {
				return err
			}
		}

		return nil
	})

	if err != nil {
		panic("failed to replay WAL")
	}

	return service
}

var _ documentgrpc.DocumentService = Service{}

func (s Service) Create(ctx context.Context, collection string, content map[string]any) (documentmodels.Document, error) {
	var documentID uuid.UUID = uuid.New()

	if err := s.log(documentmodels.OpCreate, collection, documentID, content); err != nil {
		return documentmodels.Document{}, err
	}

	s.documentStorage.Set(ctx, collection, documentID, content)
	return s.documentStorage.Get(ctx, collection, documentID)
}

func (s Service) Get(ctx context.Context, collection string, documentID string) (documentmodels.Document, error) {
	return s.documentStorage.Get(ctx, collection, uuid.MustParse(documentID))
}

func (s Service) List(ctx context.Context, collection string) ([]documentmodels.Document, error) {
	return s.documentStorage.List(ctx, collection)
}

func (s Service) Delete(ctx context.Context, collection string, documentID string) error {
	uuidID := uuid.MustParse(documentID)

	if err := s.log(documentmodels.OpDelete, collection, uuidID, nil); err != nil {
		return err
	}

	return s.documentStorage.Delete(ctx, collection, uuidID)
}

func (s Service) Update(ctx context.Context, collection string, documentID string, content map[string]any) (documentmodels.Document, error) {
	uuidID := uuid.MustParse(documentID)

	if err := s.log(documentmodels.OpUpdate, collection, uuidID, content); err != nil {
		return documentmodels.Document{}, err
	}

	return s.documentStorage.Replace(ctx, collection, uuidID, content)
}

func (s *Service) log(op documentmodels.OperationType, collection string, documentID uuid.UUID, content map[string]any) error {
	record := walmodels.Record{
		Op:         op,
		Timestamp:  time.Now(),
		Collection: collection,
		DocumentID: documentID,
		Content:    content,
	}
	data, err := json.Marshal(record)
	if err != nil {
		return err
	}
	return s.walStorage.Write(data)
}
