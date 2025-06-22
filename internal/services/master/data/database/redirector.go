package databaserdr

import (
	"context"
	"errors"

	databaseapi "github.com/10Narratives/distgo-db/internal/grpc/master/data/database"
	clustermodels "github.com/10Narratives/distgo-db/internal/models/master/cluster"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClusterStorage interface {
	Worker(databaseName string) (clustermodels.Worker, error)
}

type Redirector struct {
	clusterStorage ClusterStorage
	clients        map[string]dbv1.DatabaseServiceClient
}

var _ databaseapi.DatabaseRedirector = &Redirector{}

func New(clusterStorage ClusterStorage) *Redirector {
	return &Redirector{
		clusterStorage: clusterStorage,
		clients:        map[string]dbv1.DatabaseServiceClient{},
	}
}

func (r *Redirector) CreateDatabase(ctx context.Context, req *dbv1.CreateDatabaseRequest) (*dbv1.Database, error) {
	worker, err := r.clusterStorage.Worker(req.DatabaseId)
	if err != nil {
		return nil, errors.New("cannot create database on unknown node")
	}

	if _, contains := r.clients[worker.Database]; !contains {
		r.clients[worker.Database] = dbv1.NewDatabaseServiceClient(worker.Conn)
	}

	return r.clients[worker.Database].CreateDatabase(ctx, req)
}

func (r *Redirector) DeleteDatabase(ctx context.Context, req *dbv1.DeleteDatabaseRequest) (*emptypb.Empty, error) {
	key := databasemodels.NewKey(req.GetName())

	worker, err := r.clusterStorage.Worker(key.Database)
	if err != nil {
		return nil, errors.New("cannot create database on unknown node")
	}

	if _, contains := r.clients[worker.Database]; !contains {
		r.clients[worker.Database] = dbv1.NewDatabaseServiceClient(worker.Conn)
	}

	return r.clients[worker.Database].DeleteDatabase(ctx, req)
}

func (r *Redirector) GetDatabase(ctx context.Context, req *dbv1.GetDatabaseRequest) (*dbv1.Database, error) {
	key := databasemodels.NewKey(req.GetName())

	worker, err := r.clusterStorage.Worker(key.Database)
	if err != nil {
		return nil, errors.New("cannot create database on unknown node")
	}

	if _, contains := r.clients[worker.Database]; !contains {
		r.clients[worker.Database] = dbv1.NewDatabaseServiceClient(worker.Conn)
	}

	return r.clients[worker.Database].GetDatabase(ctx, req)
}

func (r *Redirector) ListDatabases(ctx context.Context, req *dbv1.ListDatabasesRequest) (*dbv1.ListDatabasesResponse, error) {
	panic("unimplemented")
}

func (r *Redirector) UpdateDatabase(ctx context.Context, req *dbv1.UpdateDatabaseRequest) (*dbv1.Database, error) {
	key := databasemodels.NewKey(req.GetDatabase().GetName())

	worker, err := r.clusterStorage.Worker(key.Database)
	if err != nil {
		return nil, errors.New("cannot create database on unknown node")
	}

	if _, contains := r.clients[worker.Database]; !contains {
		r.clients[worker.Database] = dbv1.NewDatabaseServiceClient(worker.Conn)
	}

	return r.clients[worker.Database].UpdateDatabase(ctx, req)
}
