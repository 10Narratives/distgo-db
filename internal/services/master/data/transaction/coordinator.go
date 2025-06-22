package transactioncdr

import (
	"context"
	"errors"

	transactionapi "github.com/10Narratives/distgo-db/internal/grpc/master/data/transaction"
	clustermodels "github.com/10Narratives/distgo-db/internal/models/master/cluster"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClusterStorage interface {
	Worker(databaseName string) (clustermodels.Worker, error)
}

type Coordinator struct {
	clusterStorage ClusterStorage
	knownClients   map[string]dbv1.TransactionServiceClient
}

func New(clusterStorage ClusterStorage) *Coordinator {
	return &Coordinator{
		clusterStorage: clusterStorage,
		knownClients:   make(map[string]dbv1.TransactionServiceClient),
	}
}

var _ transactionapi.TransactionCoordinator = &Coordinator{}

func (c *Coordinator) getClient(databaseName string) (dbv1.TransactionServiceClient, error) {
	if client, exists := c.knownClients[databaseName]; exists {
		return client, nil
	}

	worker, err := c.clusterStorage.Worker(databaseName)
	if err != nil {
		return nil, errors.New("unknown worker for database: " + databaseName)
	}

	client := dbv1.NewTransactionServiceClient(worker.Conn)
	c.knownClients[databaseName] = client
	return client, nil
}

func (c *Coordinator) Begin(ctx context.Context, req *dbv1.BeginRequest) (*dbv1.BeginResponse, error) {
	client, err := c.getClient(req.GetDatabaseName())
	if err != nil {
		return nil, err
	}
	return client.Begin(ctx, req)
}

func (c *Coordinator) Commit(ctx context.Context, req *dbv1.CommitRequest) (*emptypb.Empty, error) {
	client, err := c.getClient(req.GetDatabaseName())
	if err != nil {
		return nil, err
	}
	return client.Commit(ctx, req)
}

func (c *Coordinator) Execute(ctx context.Context, req *dbv1.ExecuteRequest) (*emptypb.Empty, error) {
	client, err := c.getClient(req.GetDatabaseName())
	if err != nil {
		return nil, err
	}
	return client.Execute(ctx, req)
}

func (c *Coordinator) Rollback(ctx context.Context, req *dbv1.RollbackRequest) (*emptypb.Empty, error) {
	client, err := c.getClient(req.GetDatabaseName())
	if err != nil {
		return nil, err
	}
	return client.Rollback(ctx, req)
}
