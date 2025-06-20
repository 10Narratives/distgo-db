package clustersrv

import (
	"context"
	"errors"

	clusterapi "github.com/10Narratives/distgo-db/internal/grpc/master/cluster"
	clustermodels "github.com/10Narratives/distgo-db/internal/models/master/cluster"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type ClusterStorage interface {
	DeleteWorker(workerID string) error
	CreateWorker(databaseName string, conn *grpc.ClientConn) (clustermodels.Worker, error)
	Worker(databaseName string) (clustermodels.Worker, error)
	Workers() []clustermodels.Worker
}

type Service struct {
	clusterStorage ClusterStorage
}

func New(clusterStorage ClusterStorage) *Service {
	return &Service{
		clusterStorage: clusterStorage,
	}
}

var _ clusterapi.ClusterService = &Service{}

func (s *Service) Register(ctx context.Context, databaseName string, address string) (string, error) {
	_, err := s.clusterStorage.Worker(databaseName)
	if err == nil {
		return "", errors.New("cannot register worker: database name already in use")
	}

	conn, err := grpc.NewClient(address, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", errors.New("failed to connect worker: " + err.Error())
	}

	worker, err := s.clusterStorage.CreateWorker(databaseName, conn)
	if err != nil {
		return "", errors.New("failed to create worker: " + err.Error())
	}

	return worker.ID.String(), nil
}

func (s *Service) Unregister(ctx context.Context, workerID string) error {
	return s.clusterStorage.DeleteWorker(workerID)
}

func (s *Service) Workers(ctx context.Context) []clustermodels.Worker {
	return s.clusterStorage.Workers()
}
