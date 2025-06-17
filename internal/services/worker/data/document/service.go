package documentsrv

import (
	"context"
	"errors"
	"time"

	documentgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/document"
	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
	commonsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/common"
)

//go:generate mockery --name DocumentStorage --output ./mocks/
type DocumentStorage interface {
	Document(ctx context.Context, key documentmodels.Key) (documentmodels.Document, error)
	Documents(ctx context.Context, parent collectionmodels.Key) []documentmodels.Document
	CreateDocument(ctx context.Context, key documentmodels.Key, value string) (documentmodels.Document, error)
	DeleteDocument(ctx context.Context, key documentmodels.Key) error
	UpdateDocument(ctx context.Context, document documentmodels.Document) error
}

type Service struct {
	documentStore DocumentStorage
	walStorage    commonsrv.WALStorage
}

var _ documentgrpc.DocumentService = &Service{}

func New(storage DocumentStorage, walStorage commonsrv.WALStorage) *Service {
	return &Service{
		documentStore: storage,
		walStorage:    walStorage,
	}
}

func (s *Service) Document(ctx context.Context, name string) (documentmodels.Document, error) {
	key := documentmodels.NewKey(name)
	return s.documentStore.Document(ctx, key)
}

func (s *Service) Documents(ctx context.Context, parent string, size int32, token string) ([]documentmodels.Document, string, error) {
	parentKey := collectionmodels.NewKey(parent)
	allDocs := s.documentStore.Documents(ctx, parentKey)

	if len(allDocs) == 0 {
		return []documentmodels.Document{}, "", nil
	}

	startIndex := 0
	if token != "" {
		for i, doc := range allDocs {
			if doc.Name == token {
				startIndex = i + 1
				break
			}
		}
	}

	endIndex := startIndex + int(size)
	if endIndex > len(allDocs) {
		endIndex = len(allDocs)
	}

	page := allDocs[startIndex:endIndex]

	var nextPageToken string
	if endIndex < len(allDocs) {
		nextPageToken = allDocs[endIndex].Name
	}

	return page, nextPageToken, nil
}

func (s *Service) CreateDocument(ctx context.Context, parent string, documentID string, value string) (documentmodels.Document, error) {
	key := documentmodels.NewKey(parent + "/documents/" + documentID)

	doc, err := s.documentStore.CreateDocument(ctx, key, value)
	if err != nil {
		return documentmodels.Document{}, err
	}

	entry := walmodels.WALEntry{
		ID:        doc.Name,
		Target:    "document",
		Type:      commonmodels.MutationTypeCreate,
		NewValue:  string(doc.Value),
		Timestamp: time.Now(),
	}
	if err := s.walStorage.LogEntry(ctx, entry); err != nil {
		return documentmodels.Document{}, errors.New("failed to log WAL entry: " + err.Error())
	}

	return doc, nil
}

func (s *Service) DeleteDocument(ctx context.Context, name string) error {
	key := documentmodels.NewKey(name)

	doc, err := s.documentStore.Document(ctx, key)
	if err != nil {
		return err
	}

	if err := s.documentStore.DeleteDocument(ctx, key); err != nil {
		return err
	}

	entry := walmodels.WALEntry{
		ID:        doc.Name,
		Target:    "document",
		Type:      commonmodels.MutationTypeDelete,
		OldValue:  string(doc.Value),
		Timestamp: time.Now(),
	}
	if err := s.walStorage.LogEntry(ctx, entry); err != nil {
		return errors.New("failed to log WAL entry: " + err.Error())
	}

	return nil
}

func (s *Service) UpdateDocument(ctx context.Context, document documentmodels.Document, paths []string) (documentmodels.Document, error) {
	for _, path := range paths {
		switch path {
		case "value":
			existingDoc, err := s.documentStore.Document(ctx, documentmodels.NewKey(document.Name))
			if err != nil {
				return documentmodels.Document{}, err
			}

			if err := s.documentStore.UpdateDocument(ctx, document); err != nil {
				return documentmodels.Document{}, err
			}

			entry := walmodels.WALEntry{
				ID:        document.Name,
				Target:    "document",
				Type:      commonmodels.MutationTypeUpdate,
				OldValue:  string(existingDoc.Value),
				NewValue:  string(document.Value),
				Timestamp: time.Now(),
			}
			if err := s.walStorage.LogEntry(ctx, entry); err != nil {
				return documentmodels.Document{}, errors.New("failed to log WAL entry: " + err.Error())
			}

			return document, nil
		default:
			return documentmodels.Document{}, errors.New("unknown field: " + path)
		}
	}

	return document, nil
}
