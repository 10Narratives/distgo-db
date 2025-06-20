package transactionstorage

import (
	"context"
	"fmt"
	"sync"

	transactionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/transaction"
	transactionsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/transaction"
)

type Storage struct {
	mu           sync.Mutex
	transactions map[string][]transactionmodels.TransactionEntry
}

var _ transactionsrv.TransactionStorage = &Storage{}

func New() *Storage {
	return &Storage{
		transactions: make(map[string][]transactionmodels.TransactionEntry),
	}
}

func (s *Storage) Activate(txID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.transactions[txID]; exists {
		return fmt.Errorf("transaction with ID %s already exists", txID)
	}

	s.transactions[txID] = []transactionmodels.TransactionEntry{}
	return nil
}

func (s *Storage) Append(ctx context.Context, entry transactionmodels.TransactionEntry) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.transactions[entry.TransactionID]; !exists {
		return fmt.Errorf("transaction with ID %s does not exist", entry.TransactionID)
	}

	s.transactions[entry.TransactionID] = append(s.transactions[entry.TransactionID], entry)
	return nil
}

func (s *Storage) GetAllByTxID(ctx context.Context, txID string) ([]transactionmodels.TransactionEntry, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	entries, exists := s.transactions[txID]
	if !exists {
		return nil, fmt.Errorf("transaction with ID %s does not exist", txID)
	}

	return entries, nil
}

func (s *Storage) Remove(txID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	if _, exists := s.transactions[txID]; !exists {
		return fmt.Errorf("transaction with ID %s does not exist", txID)
	}

	delete(s.transactions, txID)
	return nil
}
