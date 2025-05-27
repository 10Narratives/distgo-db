package documentgrpc

import (
	"context"

	documentv1 "github.com/10Narratives/distgo-db/pkg/proto/distgodb/worker/document/v1"
	"google.golang.org/grpc"
)

type DocumentService interface {
	ListDocuments(ctx context.Context, collection string) ([]documentv1.Document, error)
	GetDocument(ctx context.Context, collection string, documentID string) (documentv1.Document, error)
	CreateDocument(ctx context.Context, collection string, data map[string]any) (documentv1.Document, error)
	UpdateDocument(ctx context.Context, collection string, update map[string]any) (documentv1.Document, error)
	DeleteDocument(ctx context.Context, collection, documentID string) (bool, error)
}

type serverAPI struct {
	documentv1.UnimplementedDocumentServiceServer
	documentSrv DocumentService
}

func Register(gRPCServer *grpc.Server, documentSrv DocumentService) {
	documentv1.RegisterDocumentServiceServer(gRPCServer, &serverAPI{documentSrv: documentSrv})
}

func (s *serverAPI) ListDocuments(ctx context.Context, req *documentv1.ListDocumentsRequest) (*documentv1.ListDocumentsResponse, error) {
	panic("implement me")
	// return nil, nil
}

func (s *serverAPI) GetDocument(ctx context.Context, req *documentv1.GetDocumentRequest) (*documentv1.GetDocumentResponse, error) {
	panic("implement me")
	// return nil, nil
}

func (s *serverAPI) CreateDocument(ctx context.Context, req *documentv1.CreateDocumentRequest) (*documentv1.CreateDocumentResponse, error) {
	panic("implement me")
	// return nil, nil
}

func (s *serverAPI) UpdateDocument(ctx context.Context, req *documentv1.UpdateDocumentRequest) (*documentv1.UpdateDocumentResponse, error) {
	panic("implement me")
	// return nil, nil
}

func (s *serverAPI) DeleteDocument(ctx context.Context, req *documentv1.DeleteDocumentRequest) (*documentv1.DeleteDocumentResponse, error) {
	panic("implement me")
	// return nil, nil
}
