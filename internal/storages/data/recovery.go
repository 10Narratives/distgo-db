package datastorage

import (
	"bufio"
	"context"
	"encoding/json"
	"fmt"
	"os"

	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
)

func (s *Storage) RecoverFromFile(ctx context.Context, filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("failed to open file: %w", err)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Bytes()

		var entry walmodels.WALEntry
		if err := json.Unmarshal(line, &entry); err != nil {
			return fmt.Errorf("failed to unmarshal WAL entry: %w", err)
		}

		switch entry.Entity {
		case commonmodels.EntityTypeDatabase:
			if err := s.ApplyToDatabase(entry); err != nil {
				return fmt.Errorf("failed to recover database: %w", err)
			}
		case commonmodels.EntityTypeCollection:
			if err := s.ApplyToCollection(entry); err != nil {
				return fmt.Errorf("failed to recover collection: %w", err)
			}
		case commonmodels.EntityTypeDocument:
			if err := s.ApplyToDocument(entry); err != nil {
				return fmt.Errorf("failed to recover document: %w", err)
			}
		default:
			return fmt.Errorf("unknown entity type: %d", entry.Entity)
		}
	}

	if err := scanner.Err(); err != nil {
		return fmt.Errorf("error while reading file: %w", err)
	}

	return nil
}

func (s *Storage) ApplyToDatabase(entry walmodels.WALEntry) error {
	var payload walmodels.DatabasePayload
	if err := json.Unmarshal(entry.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal database payload: %w", err)
	}

	switch entry.Mutation {
	case commonmodels.MutationTypeCreate:
		if payload.Database == nil {
			return fmt.Errorf("database payload is nil for create mutation")
		}
		s.databases.Store(payload.Key, *payload.Database)
	case commonmodels.MutationTypeDelete:
		s.databases.Delete(payload.Key)
	default:
		return fmt.Errorf("unsupported mutation type for database: %d", entry.Mutation)
	}

	return nil
}

func (s *Storage) ApplyToCollection(entry walmodels.WALEntry) error {
	var payload walmodels.CollectionPayload
	if err := json.Unmarshal(entry.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal collection payload: %w", err)
	}

	switch entry.Mutation {
	case commonmodels.MutationTypeCreate:
		if payload.Collection == nil {
			return fmt.Errorf("collection payload is nil for create mutation")
		}
		s.collections.Store(payload.Key, *payload.Collection)
	case commonmodels.MutationTypeDelete:
		s.collections.Delete(payload.Key)
	default:
		return fmt.Errorf("unsupported mutation type for collection: %d", entry.Mutation)
	}

	return nil
}

func (s *Storage) ApplyToDocument(entry walmodels.WALEntry) error {
	var payload walmodels.DocumentPayload
	if err := json.Unmarshal(entry.Payload, &payload); err != nil {
		return fmt.Errorf("failed to unmarshal document payload: %w", err)
	}

	switch entry.Mutation {
	case commonmodels.MutationTypeCreate, commonmodels.MutationTypeUpdate:
		if payload.Document == nil {
			return fmt.Errorf("document payload is nil for create/update mutation")
		}
		s.documents.Store(payload.Key, *payload.Document)
	case commonmodels.MutationTypeDelete:
		s.documents.Delete(payload.Key)
	default:
		return fmt.Errorf("unsupported mutation type for document: %d", entry.Mutation)
	}

	return nil
}
