package transactionstorage

// import (
// 	"sync"

// 	commonmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/common"
// 	transactionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/transaction"
// 	transactionsrv "github.com/10Narratives/distgo-db/internal/services/worker/data/transaction"
// )

// type Storage struct {
// 	mu           sync.Mutex
// 	transactions map[string]*transactionmodels.Transaction
// }

// var _ transactionsrv.TransactionStorage = &Storage{}

// func New() *Storage {
// 	return &Storage{
// 		transactions: make(map[string]*transactionmodels.Transaction),
// 	}
// }

// func (s *Storage) Add(txID string, operations []commonmodels.Operation) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	if _, exists := s.transactions[txID]; !exists {
// 		s.transactions[txID] = &transactionmodels.Transaction{
// 			ID: txID,
// 		}
// 	}
// 	s.transactions[txID].Operations = append(s.transactions[txID].Operations, operations...)
// }

// func (s *Storage) Delete(txID string) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	delete(s.transactions, txID)
// }

// func (s *Storage) Get(txID string) (*transactionmodels.Transaction, bool) {
// 	s.mu.Lock()
// 	defer s.mu.Unlock()

// 	tx, exists := s.transactions[txID]
// 	return tx, exists
// }
