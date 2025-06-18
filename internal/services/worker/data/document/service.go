package documentsrv

import (
	"context"
	"errors"
	"time"

	documentgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/document"
	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
)

//go:generate mockery --name DocumentStorage --output ./mocks/
type DocumentStorage interface {
	Document(ctx context.Context, key documentmodels.Key) (documentmodels.Document, error)
	Documents(ctx context.Context, parent collectionmodels.Key) ([]documentmodels.Document, error)
	CreateDocument(ctx context.Context, key documentmodels.Key, value string) (documentmodels.Document, error)
	DeleteDocument(ctx context.Context, key documentmodels.Key) error
	UpdateDocument(ctx context.Context, document documentmodels.Document) error
}

//go:generate mockery --name WALService --output ./mocks/
type WALService interface {
	CreateDocumentEntry(ctx context.Context, mutation commonmodels.MutationType, key documentmodels.Key, doc *documentmodels.Document) error
}

type Service struct {
	documentStorage DocumentStorage
	walService      WALService
}

var _ documentgrpc.DocumentService = &Service{}

func New(documentStorage DocumentStorage, walService WALService) *Service {
	return &Service{
		documentStorage: documentStorage,
		walService:      walService,
	}
}

func (s *Service) Document(ctx context.Context, name string) (documentmodels.Document, error) {
	key := documentmodels.NewKey(name)
	return s.documentStorage.Document(ctx, key)
}

func (s *Service) Documents(ctx context.Context, parent string, size int32, token string) ([]documentmodels.Document, string, error) {
	parentKey := collectionmodels.NewKey(parent)
	allDocs, err := s.documentStorage.Documents(ctx, parentKey)
	if err != nil {
		return []documentmodels.Document{}, "", err
	}

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

	newDoc := documentmodels.Document{
		Name:      parent + "/documents/" + documentID,
		ID:        documentID,
		Value:     []byte(value),
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
		Parent:    parent,
	}

	if err := s.walService.CreateDocumentEntry(ctx, commonmodels.MutationTypeCreate, key, &newDoc); err != nil {
		return documentmodels.Document{}, errors.New("failed to create WAL entry: " + err.Error())
	}

	doc, err := s.documentStorage.CreateDocument(ctx, key, value)
	if err != nil {
		return documentmodels.Document{}, err
	}

	return doc, nil
}

func (s *Service) DeleteDocument(ctx context.Context, name string) error {
	key := documentmodels.NewKey(name)

	if err := s.walService.CreateDocumentEntry(ctx, commonmodels.MutationTypeDelete, key, nil); err != nil {
		return errors.New("failed to create WAL entry: " + err.Error())
	}

	if err := s.documentStorage.DeleteDocument(ctx, key); err != nil {
		return err
	}

	return nil
}

func (s *Service) UpdateDocument(ctx context.Context, document documentmodels.Document, paths []string) (documentmodels.Document, error) {
	key := documentmodels.NewKey(document.Name)

	for _, path := range paths {
		switch path {
		case "value":
			updatedDoc := document
			updatedDoc.UpdatedAt = time.Now()

			if err := s.walService.CreateDocumentEntry(ctx, commonmodels.MutationTypeUpdate, key, &updatedDoc); err != nil {
				return documentmodels.Document{}, errors.New("failed to create WAL entry: " + err.Error())
			}

			if err := s.documentStorage.UpdateDocument(ctx, updatedDoc); err != nil {
				return documentmodels.Document{}, err
			}

			return updatedDoc, nil
		default:
			return documentmodels.Document{}, errors.New("unknown field: " + path)
		}
	}

	return document, nil
}
