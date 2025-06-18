package datastorage

// import (
// 	"context"
// 	"encoding/json"
// 	"fmt"

// 	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
// 	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
// 	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
// 	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
// 	walmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/wal"
// )

// // Recover applies WAL entries to reconstruct the storage state.
// func (s *Storage) Recover(ctx context.Context, walEntries []walmodels.WALEntry) error {
// 	for _, entry := range walEntries {
// 		switch entry.Target {
// 		case "database":
// 			if err := s.applyDatabaseEntry(entry); err != nil {
// 				return fmt.Errorf("failed to apply database WAL entry: %w", err)
// 			}
// 		case "collection":
// 			if err := s.applyCollectionEntry(entry); err != nil {
// 				return fmt.Errorf("failed to apply collection WAL entry: %w", err)
// 			}
// 		case "document":
// 			if err := s.applyDocumentEntry(entry); err != nil {
// 				return fmt.Errorf("failed to apply document WAL entry: %w", err)
// 			}
// 		default:
// 			return fmt.Errorf("unknown target in WAL entry: %s", entry.Target)
// 		}
// 	}
// 	return nil
// }

// func (s *Storage) applyDatabaseEntry(entry walmodels.WALEntry) error {
// 	key := databasemodels.NewKey(entry.ID)

// 	switch entry.Type {
// 	case commonmodels.MutationTypeCreate:
// 		var db databasemodels.Database
// 		if err := json.Unmarshal([]byte(entry.NewValue), &db); err != nil {
// 			return fmt.Errorf("failed to unmarshal database: %w", err)
// 		}
// 		s.databases.Store(key, db)
// 	case commonmodels.MutationTypeUpdate:
// 		existing, ok := s.databases.Load(key)
// 		if !ok {
// 			return fmt.Errorf("database not found for update: %s", entry.ID)
// 		}
// 		db := existing.(databasemodels.Database)
// 		if err := json.Unmarshal([]byte(entry.NewValue), &db); err != nil {
// 			return fmt.Errorf("failed to unmarshal updated database: %w", err)
// 		}
// 		s.databases.Store(key, db)
// 	case commonmodels.MutationTypeDelete:
// 		s.databases.Delete(key)
// 	default:
// 		return fmt.Errorf("unknown mutation type for database: %d", entry.Type)
// 	}
// 	return nil
// }

// func (s *Storage) applyCollectionEntry(entry walmodels.WALEntry) error {
// 	key := collectionmodels.NewKey(entry.ID)

// 	switch entry.Type {
// 	case commonmodels.MutationTypeCreate:
// 		var coll collectionmodels.Collection
// 		if err := json.Unmarshal([]byte(entry.NewValue), &coll); err != nil {
// 			return fmt.Errorf("failed to unmarshal collection: %w", err)
// 		}
// 		s.collections.Store(key, coll)
// 	case commonmodels.MutationTypeUpdate:
// 		existing, ok := s.collections.Load(key)
// 		if !ok {
// 			return fmt.Errorf("collection not found for update: %s", entry.ID)
// 		}
// 		coll := existing.(collectionmodels.Collection)
// 		if err := json.Unmarshal([]byte(entry.NewValue), &coll); err != nil {
// 			return fmt.Errorf("failed to unmarshal updated collection: %w", err)
// 		}
// 		s.collections.Store(key, coll)
// 	case commonmodels.MutationTypeDelete:
// 		s.collections.Delete(key)
// 	default:
// 		return fmt.Errorf("unknown mutation type for collection: %d", entry.Type)
// 	}
// 	return nil
// }

// func (s *Storage) applyDocumentEntry(entry walmodels.WALEntry) error {
// 	key := documentmodels.NewKey(entry.ID)

// 	switch entry.Type {
// 	case commonmodels.MutationTypeCreate:
// 		var doc documentmodels.Document
// 		if err := json.Unmarshal([]byte(entry.NewValue), &doc); err != nil {
// 			return fmt.Errorf("failed to unmarshal document: %w", err)
// 		}
// 		s.documents.Store(key, doc)
// 	case commonmodels.MutationTypeUpdate:
// 		existing, ok := s.documents.Load(key)
// 		if !ok {
// 			return fmt.Errorf("document not found for update: %s", entry.ID)
// 		}
// 		doc := existing.(documentmodels.Document)
// 		if err := json.Unmarshal([]byte(entry.NewValue), &doc); err != nil {
// 			return fmt.Errorf("failed to unmarshal updated document: %w", err)
// 		}
// 		s.documents.Store(key, doc)
// 	case commonmodels.MutationTypeDelete:
// 		s.documents.Delete(key)
// 	default:
// 		return fmt.Errorf("unknown mutation type for document: %d", entry.Type)
// 	}
// 	return nil
// }
