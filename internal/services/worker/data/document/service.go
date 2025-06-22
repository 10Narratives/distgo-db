package documentsrv

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	documentgrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/data/document"
	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
)

type TransactionService interface {
	Execute(ctx context.Context, operations []commonmodels.Operation) error
}

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
	transactionSrv  TransactionService
}

var _ documentgrpc.DocumentService = &Service{}

func New(documentStorage DocumentStorage, walService WALService, transactionService TransactionService) *Service {
	return &Service{
		documentStorage: documentStorage,
		walService:      walService,
		transactionSrv:  transactionService,
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
	collKey := collectionmodels.NewKey(parent)
	if collKey.Database == "" || collKey.Collection == "" {
		return documentmodels.Document{}, errors.New("invalid collection name")
	}

	docKey := documentmodels.Key{
		Database:   collKey.Database,
		Collection: collKey.Collection,
		Document:   documentID,
	}

	fullName := docKey.FullName()

	valueJSON, err := json.Marshal(value)
	if err != nil {
		return documentmodels.Document{}, fmt.Errorf("value serialization failed: %w", err)
	}

	op := commonmodels.Operation{
		Mutation: commonmodels.MutationTypeCreate,
		Entity:   commonmodels.EntityTypeDocument,
		Name:     fullName,
		Value:    valueJSON,
	}

	if err := s.transactionSrv.Execute(ctx, []commonmodels.Operation{op}); err != nil {
		return documentmodels.Document{}, fmt.Errorf("transaction failed: %w", err)
	}

	return s.documentStorage.Document(ctx, docKey)
}

func (s *Service) DeleteDocument(ctx context.Context, name string) error {
	key := documentmodels.NewKey(name)
	if err := key.Validate(); err != nil {
		return err
	}

	op := commonmodels.Operation{
		Mutation: commonmodels.MutationTypeDelete,
		Entity:   commonmodels.EntityTypeDocument,
		Name:     key.FullName(),
		Value:    json.RawMessage("null"),
	}

	return s.transactionSrv.Execute(ctx, []commonmodels.Operation{op})
}

func (s *Service) UpdateDocument(ctx context.Context, document documentmodels.Document, paths []string) (documentmodels.Document, error) {
	key := documentmodels.NewKey(document.Name)
	if err := key.Validate(); err != nil {
		return documentmodels.Document{}, err
	}

	valueJSON, err := json.Marshal(string(document.Value))
	if err != nil {
		return documentmodels.Document{}, fmt.Errorf("value serialization failed: %w", err)
	}

	op := commonmodels.Operation{
		Mutation: commonmodels.MutationTypeUpdate,
		Entity:   commonmodels.EntityTypeDocument,
		Name:     key.FullName(),
		Value:    valueJSON,
	}

	if err := s.transactionSrv.Execute(ctx, []commonmodels.Operation{op}); err != nil {
		return documentmodels.Document{}, fmt.Errorf("transaction failed: %w", err)
	}

	return s.documentStorage.Document(ctx, key)
}
func (s *Service) BatchUpdate(ctx context.Context, operations []commonmodels.Operation) error {

	for i, op := range operations {
		if op.Entity != commonmodels.EntityTypeDocument {
			return fmt.Errorf("unsupported entity type: %s at index %d", op.Entity, i)
		}

		key := documentmodels.NewKey(op.Name)
		if key.Database == "" || key.Collection == "" || key.Document == "" {
			return fmt.Errorf("invalid document name: %s at index %d", op.Name, i)
		}
	}

	return s.transactionSrv.Execute(ctx, operations)
}
