package clusterstorage

import (
	"errors"
	"sync"

	clustermodels "github.com/10Narratives/distgo-db/internal/models/master/cluster"
	clustersrv "github.com/10Narratives/distgo-db/internal/services/master/cluster"
	"github.com/google/uuid"
	"google.golang.org/grpc"
)

type Storage struct {
	mu      sync.RWMutex
	workers map[uuid.UUID]clustermodels.Worker
}

func New() *Storage {
	return &Storage{
		workers: make(map[uuid.UUID]clustermodels.Worker),
	}
}

var _ clustersrv.ClusterStorage = &Storage{}

func (s *Storage) CreateWorker(databaseName string, conn *grpc.ClientConn) (clustermodels.Worker, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	workerID := uuid.New()

	worker := clustermodels.Worker{
		ID:       workerID,
		Database: databaseName,
		Address:  conn.Target(),
		Conn:     conn,
	}

	s.workers[workerID] = worker
	return worker, nil
}

func (s *Storage) DeleteWorker(workerID string) error {
	s.mu.Lock()
	defer s.mu.Unlock()

	id, err := uuid.Parse(workerID)
	if err != nil {
		return errors.New("invalid worker ID format")
	}

	worker, exists := s.workers[id]
	if !exists {
		return errors.New("worker not found")
	}

	if worker.Conn != nil {
		worker.Conn.Close()
	}

	delete(s.workers, id)
	return nil
}

func (s *Storage) Worker(databaseName string) (clustermodels.Worker, error) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, worker := range s.workers {
		if worker.Database == databaseName {
			return worker, nil
		}
	}

	return clustermodels.Worker{}, errors.New("worker not found")
}

func (s *Storage) Workers() []clustermodels.Worker {
	s.mu.RLock()
	defer s.mu.RUnlock()

	workers := make([]clustermodels.Worker, 0, len(s.workers))
	for _, worker := range s.workers {
		workers = append(workers, worker)
	}
	return workers
}
