package storagegrpc

import (
	"context"

	"github.com/10Narratives/distgo-db/pkg/proto/worker"
	"google.golang.org/grpc"
)

type Storage interface {
	Create(ctx context.Context, collection, documentID string)
	Get(ctx context.Context, collection, documentID string)
	List(ctx context.Context, collection string)
	Update(ctx context.Context, collection, documentID string, update map[string]any)
	Delete(ctx context.Context, collection, documentID string)
}

type serverAPI struct {
	worker.UnimplementedStorageServer
	store Storage
}

func Register(gRPCServer *grpc.Server, store Storage) {
	worker.RegisterStorageServer(gRPCServer, &serverAPI{store: store})
}

func (s serverAPI) CreateDocument(ctx context.Context, req *worker.CreateDocumentRequest) (*worker.CreateDocumentResponse, error) {
	panic("implement me")
}

func (s serverAPI) GetDocument(ctx context.Context, req *worker.GetDocumentRequest) (*worker.GetDocumentResponse, error) {
	panic("implement me")
}

func (s serverAPI) ListDocuments(ctx context.Context, req *worker.ListDocumentsRequest) (*worker.ListDocumentsResponse, error) {
	panic("implement me")
}

func (s serverAPI) UpdateDocument(ctx context.Context, req *worker.UpdateDocumentRequest) (*worker.UpdateDocumentResponse, error) {
	panic("implement me")
}

func (s serverAPI) DeleteDocument(ctx context.Context, req *worker.DeleteDocumentRequest) (*worker.DeleteDocumentResponse, error) {
	panic("implement me")
}
