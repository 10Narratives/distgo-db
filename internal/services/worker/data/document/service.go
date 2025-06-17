package documentsrv

import (
	"context"
	"errors"

	documentgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/document"
	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
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
}

var _ documentgrpc.DocumentService = &Service{}

func New(storage DocumentStorage) *Service {
	return &Service{
		documentStore: storage,
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
	return s.documentStore.CreateDocument(ctx, key, value)
}

func (s *Service) DeleteDocument(ctx context.Context, name string) error {
	key := documentmodels.NewKey(name)
	return s.documentStore.DeleteDocument(ctx, key)
}

func (s *Service) UpdateDocument(ctx context.Context, document documentmodels.Document, paths []string) (documentmodels.Document, error) {
	for _, path := range paths {
		switch path {
		case "value":
			err := s.documentStore.UpdateDocument(ctx, document)
			if err != nil {
				return documentmodels.Document{}, err
			}
		default:
			return documentmodels.Document{}, errors.New("unknown field: " + path)
		}
	}
	return document, nil
}
