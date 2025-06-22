package datastorage

import (
	"errors"
	"sync"
	"time"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
	"github.com/google/uuid"
)

type Transaction struct {
	mu        sync.RWMutex
	documents map[documentmodels.Key]documentmodels.Document
	created   map[documentmodels.Key]bool
	deleted   map[documentmodels.Key]bool
}

func newTransaction() *Transaction {
	return &Transaction{
		documents: make(map[documentmodels.Key]documentmodels.Document),
		created:   make(map[documentmodels.Key]bool),
		deleted:   make(map[documentmodels.Key]bool),
	}
}

func (s *Storage) BeginTx() string {
	txID := uuid.New().String()
	s.transactions.Store(txID, newTransaction())
	return txID
}

func (s *Storage) CommitTx(txID string) error {
	rawTx, ok := s.transactions.Load(txID)
	if !ok {
		return errors.New("transaction not found")
	}

	tx, ok := rawTx.(*Transaction)
	if !ok {
		return errors.New("invalid transaction type")
	}

	tx.mu.Lock()
	defer tx.mu.Unlock()

	for key, doc := range tx.documents {
		s.documents.Store(key, doc)
	}
	for key := range tx.deleted {
		s.documents.Delete(key)
	}

	s.transactions.Delete(txID)
	return nil
}

func (s *Storage) RollbackTx(txID string) error {
	s.transactions.Delete(txID)
	return nil
}

func (s *Storage) CreateDocumentInTx(txID string, key documentmodels.Key, value string) error {
	if key.Database == "" || key.Collection == "" || key.Document == "" {
		return ErrInvalidKey
	}

	rawTx, ok := s.transactions.Load(txID)
	if !ok {
		return errors.New("transaction not found")
	}

	tx, ok := rawTx.(*Transaction)
	if !ok {
		return errors.New("invalid transaction type")
	}

	tx.mu.Lock()
	defer tx.mu.Unlock()

	if tx.deleted[key] {
		return errors.New("cannot create document: marked for deletion in transaction")
	}
	if _, exists := tx.documents[key]; exists {
		return ErrDatabaseAlreadyExists
	}

	now := time.Now().UTC()
	tx.documents[key] = documentmodels.Document{
		Name:      "databases/" + key.Database + "/collections/" + key.Collection + "/documents/" + key.Document,
		Value:     []byte(value),
		CreatedAt: now,
		UpdatedAt: now,
		Parent:    "databases/" + key.Database + "/collections/" + key.Collection,
	}
	tx.created[key] = true

	s.transactions.Store(txID, tx)
	return nil
}

func (s *Storage) UpdateDocumentInTx(txID string, key documentmodels.Key, value string) error {
	if key.Database == "" || key.Collection == "" || key.Document == "" {
		return ErrInvalidKey
	}

	rawTx, ok := s.transactions.Load(txID)
	if !ok {
		return errors.New("transaction not found")
	}

	tx, ok := rawTx.(*Transaction)
	if !ok {
		return errors.New("invalid transaction type")
	}

	tx.mu.Lock()
	defer tx.mu.Unlock()

	if tx.deleted[key] {
		return errors.New("cannot update document: marked for deletion in transaction")
	}

	doc, exists := tx.documents[key]
	if !exists {
		rawDoc, ok := s.documents.Load(key)
		if !ok {
			return ErrDocumentNotFound
		}
		doc = rawDoc.(documentmodels.Document)
	}

	doc.Value = []byte(value)
	doc.UpdatedAt = time.Now().UTC()
	tx.documents[key] = doc

	s.transactions.Store(txID, tx)
	return nil
}

func (s *Storage) DeleteDocumentInTx(txID string, key documentmodels.Key) error {
	if key.Database == "" || key.Collection == "" || key.Document == "" {
		return ErrInvalidKey
	}

	rawTx, ok := s.transactions.Load(txID)
	if !ok {
		return errors.New("transaction not found")
	}

	tx, ok := rawTx.(*Transaction)
	if !ok {
		return errors.New("invalid transaction type")
	}

	tx.mu.Lock()
	defer tx.mu.Unlock()
	switch {
	case tx.created[key]:
		delete(tx.documents, key)
		delete(tx.created, key)

	case tx.deleted[key]:
		return nil

	default:
		if _, ok := tx.documents[key]; !ok {
			if _, exists := s.documents.Load(key); !exists {
				return ErrDocumentNotFound
			}
		}

		delete(tx.documents, key)
		tx.deleted[key] = true
	}

	s.transactions.Store(txID, tx)
	return nil
}

func (s *Storage) DocumentInTx(txID string, key documentmodels.Key) (documentmodels.Document, error) {
	rawTx, ok := s.transactions.Load(txID)
	if !ok {
		return documentmodels.Document{}, errors.New("transaction not found")
	}

	tx, ok := rawTx.(*Transaction)
	if !ok {
		return documentmodels.Document{}, errors.New("invalid transaction type")
	}

	tx.mu.RLock()
	defer tx.mu.RUnlock()

	if tx.deleted[key] {
		return documentmodels.Document{}, errors.New("document deleted in transaction")
	}

	if doc, exists := tx.documents[key]; exists {
		return doc, nil
	}

	rawDoc, ok := s.documents.Load(key)
	if !ok {
		return documentmodels.Document{}, ErrDocumentNotFound
	}

	return rawDoc.(documentmodels.Document), nil
}
