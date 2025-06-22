package documentcdr

import (
	"context"
	"errors"

	documentapi "github.com/10Narratives/distgo-db/internal/grpc/master/data/document"
	clustermodels "github.com/10Narratives/distgo-db/internal/models/master/cluster"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/data/database"
	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClusterStorage interface {
	Worker(databaseName string) (clustermodels.Worker, error)
}

type Coordinator struct {
	clusterStorage ClusterStorage
	knownClients   map[string]dbv1.DocumentServiceClient
}

func New(clusterStorage ClusterStorage) *Coordinator {
	return &Coordinator{
		clusterStorage: clusterStorage,
		knownClients:   make(map[string]dbv1.DocumentServiceClient),
	}
}

var _ documentapi.DocumentCoordinator = &Coordinator{}

func (c *Coordinator) getClient(databaseName string) (dbv1.DocumentServiceClient, error) {
	if client, exists := c.knownClients[databaseName]; exists {
		return client, nil
	}

	worker, err := c.clusterStorage.Worker(databaseName)
	if err != nil {
		return nil, errors.New("unknown worker for database: " + databaseName)
	}

	client := dbv1.NewDocumentServiceClient(worker.Conn)
	c.knownClients[databaseName] = client
	return client, nil
}

func (c *Coordinator) CreateDocument(ctx context.Context, req *dbv1.CreateDocumentRequest) (*dbv1.Document, error) {
	key := databasemodels.NewKey(req.GetParent())
	client, err := c.getClient(key.Database)
	if err != nil {
		return nil, errors.New("cannot create document: " + err.Error())
	}
	return client.CreateDocument(ctx, req)
}

func (c *Coordinator) DeleteDocument(ctx context.Context, req *dbv1.DeleteDocumentRequest) (*emptypb.Empty, error) {
	key := documentmodels.NewKey(req.GetName())
	client, err := c.getClient(key.Database)
	if err != nil {
		return nil, errors.New("cannot delete document: " + err.Error())
	}
	return client.DeleteDocument(ctx, req)
}

func (c *Coordinator) GetDocument(ctx context.Context, req *dbv1.GetDocumentRequest) (*dbv1.Document, error) {
	key := documentmodels.NewKey(req.GetName())
	client, err := c.getClient(key.Database)
	if err != nil {
		return nil, errors.New("cannot get document: " + err.Error())
	}
	return client.GetDocument(ctx, req)
}

func (c *Coordinator) ListDocuments(ctx context.Context, req *dbv1.ListDocumentsRequest) (*dbv1.ListDocumentsResponse, error) {
	key := databasemodels.NewKey(req.GetParent())
	client, err := c.getClient(key.Database)
	if err != nil {
		return nil, errors.New("cannot list documents: " + err.Error())
	}
	return client.ListDocuments(ctx, req)
}

func (c *Coordinator) UpdateDocument(ctx context.Context, req *dbv1.UpdateDocumentRequest) (*dbv1.Document, error) {
	key := documentmodels.NewKey(req.GetDocument().GetName())
	client, err := c.getClient(key.Database)
	if err != nil {
		return nil, errors.New("cannot update document: " + err.Error())
	}
	return client.UpdateDocument(ctx, req)
}
