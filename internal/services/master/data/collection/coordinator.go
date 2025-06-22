package collectioncdr

import (
	"context"
	"errors"

	collectionapi "github.com/10Narratives/distgo-db/internal/grpc/master/data/collection"
	clustermodels "github.com/10Narratives/distgo-db/internal/models/master/cluster"
	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClusterStorage interface {
	Worker(databaseName string) (clustermodels.Worker, error)
}

type Coordinator struct {
	clusterStorage ClusterStorage
	knownClients   map[string]dbv1.CollectionServiceClient
}

func New(clusterStorage ClusterStorage) *Coordinator {
	return &Coordinator{
		clusterStorage: clusterStorage,
		knownClients:   make(map[string]dbv1.CollectionServiceClient),
	}
}

var _ collectionapi.CollectionCoordinator = &Coordinator{}

func (c *Coordinator) getClient(databaseName string) (dbv1.CollectionServiceClient, error) {
	if client, exists := c.knownClients[databaseName]; exists {
		return client, nil
	}

	worker, err := c.clusterStorage.Worker(databaseName)
	if err != nil {
		return nil, errors.New("unknown worker for database: " + databaseName)
	}

	client := dbv1.NewCollectionServiceClient(worker.Conn)
	c.knownClients[databaseName] = client
	return client, nil
}

func (c *Coordinator) CreateCollection(ctx context.Context, req *dbv1.CreateCollectionRequest) (*dbv1.Collection, error) {
	key := databasemodels.NewKey(req.GetParent())
	client, err := c.getClient(key.Database)
	if err != nil {
		return nil, errors.New("cannot create collection: " + err.Error())
	}
	return client.CreateCollection(ctx, req)
}

func (c *Coordinator) DeleteCollection(ctx context.Context, req *dbv1.DeleteCollectionRequest) (*emptypb.Empty, error) {
	key := collectionmodels.NewKey(req.GetName())
	client, err := c.getClient(key.Database)
	if err != nil {
		return nil, errors.New("cannot delete collection: " + err.Error())
	}
	return client.DeleteCollection(ctx, req)
}

func (c *Coordinator) GetCollection(ctx context.Context, req *dbv1.GetCollectionRequest) (*dbv1.Collection, error) {
	key := collectionmodels.NewKey(req.GetName())
	client, err := c.getClient(key.Database)
	if err != nil {
		return nil, errors.New("cannot get collection: " + err.Error())
	}
	return client.GetCollection(ctx, req)
}

func (c *Coordinator) ListCollections(ctx context.Context, req *dbv1.ListCollectionsRequest) (*dbv1.ListCollectionsResponse, error) {
	key := databasemodels.NewKey(req.GetParent())
	client, err := c.getClient(key.Database)
	if err != nil {
		return nil, errors.New("cannot list collections: " + err.Error())
	}
	return client.ListCollections(ctx, req)
}

func (c *Coordinator) UpdateCollection(ctx context.Context, req *dbv1.UpdateCollectionRequest) (*dbv1.Collection, error) {
	key := collectionmodels.NewKey(req.GetCollection().GetName())
	client, err := c.getClient(key.Database)
	if err != nil {
		return nil, errors.New("cannot update collection: " + err.Error())
	}
	return client.UpdateCollection(ctx, req)
}
