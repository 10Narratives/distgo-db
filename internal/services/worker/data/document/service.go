package documentsrv

import (
	"context"
	"errors"

	documentgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/document"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
)

//go:generate mockery --name DocumentStorage --output ./mocks/
type DocumentStorage interface {
	Document(ctx context.Context, name string) (documentmodels.Document, error)
	Documents(ctx context.Context, parent string) []documentmodels.Document
	CreateDocument(ctx context.Context, name, value string) (documentmodels.Document, error)
	DeleteDocument(ctx context.Context, name string) error
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
	return s.documentStore.Document(ctx, name)
}

func (s *Service) Documents(ctx context.Context, parent string, size int32, token string) ([]documentmodels.Document, string, error) {
	allDocs := s.documentStore.Documents(ctx, parent)

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
	name := parent + "/documents/" + documentID
	return s.documentStore.CreateDocument(ctx, name, value)
}

func (s *Service) DeleteDocument(ctx context.Context, name string) error {
	return s.documentStore.DeleteDocument(ctx, name)
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
