package transactionsrv

import (
	"context"
	"encoding/json"
	"fmt"

	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
	datastorage "github.com/10Narratives/distgo-db/internal/storages/data"
)

type WALService interface {
	LogTransaction(ctx context.Context, operations []commonmodels.Operation) error
}

type Service struct {
	storage    *datastorage.Storage
	walService WALService
}

func New(storage *datastorage.Storage, walService WALService) *Service {
	return &Service{
		storage:    storage,
		walService: walService,
	}
}

func (s *Service) Execute(ctx context.Context, operations []commonmodels.Operation) error {
	txID := s.storage.BeginTx()

	defer func() {
		if r := recover(); r != nil {
			s.storage.RollbackTx(txID)
			panic(r)
		}
	}()

	for i, op := range operations {
		if op.Entity != commonmodels.EntityTypeDocument {
			s.storage.RollbackTx(txID)
			return fmt.Errorf("unsupported entity type: %s at index %d", op.Entity, i)
		}

		key := documentmodels.NewKey(op.Name)
		if err := key.Validate(); err != nil {
			s.storage.RollbackTx(txID)
			return fmt.Errorf("invalid document name: %s at index %d: %w",
				op.Name, i, err)
		}

		var err error
		switch op.Mutation {
		case commonmodels.MutationTypeCreate, commonmodels.MutationTypeUpdate:
			// Для операций создания и обновления передаем сырые JSON-данные
			err = s.processDocumentMutation(txID, op.Mutation, key, op.Value, i)

		case commonmodels.MutationTypeDelete:
			err = s.storage.DeleteDocumentInTx(txID, key)

		default:
			s.storage.RollbackTx(txID)
			return fmt.Errorf("unknown mutation type: %s at index %d", op.Mutation, i)
		}

		if err != nil {
			s.storage.RollbackTx(txID)
			return fmt.Errorf("operation failed [%d] mutation:%s entity:%s name:%s: %w",
				i, op.Mutation, op.Entity, op.Name, err)
		}
	}

	if err := s.walService.LogTransaction(ctx, operations); err != nil {
		s.storage.RollbackTx(txID)
		return fmt.Errorf("WAL write failed: %w", err)
	}

	if err := s.storage.CommitTx(txID); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}

func (s *Service) processDocumentMutation(
	txID string,
	mutation commonmodels.MutationType,
	key documentmodels.Key,
	value json.RawMessage,
	index int,
) error {
	rawValue := string(value)

	switch mutation {
	case commonmodels.MutationTypeCreate:
		return s.storage.CreateDocumentInTx(txID, key, rawValue)
	case commonmodels.MutationTypeUpdate:
		return s.storage.UpdateDocumentInTx(txID, key, rawValue)
	default:
		return fmt.Errorf("unsupported mutation type for document: %s at index %d", mutation, index)
	}
}
