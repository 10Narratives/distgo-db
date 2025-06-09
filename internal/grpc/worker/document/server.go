package documentgrpc

import (
	"context"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

//go:generate mockery --name DocumentService --output ./mocks/
type DocumentService interface {
	Create(ctx context.Context, collection string, content map[string]any) (documentmodels.Document, error)
	Get(ctx context.Context, collection, documentID string) (documentmodels.Document, error)
	List(ctx context.Context, collection string) ([]documentmodels.Document, error)
	Delete(ctx context.Context, collection, documentID string) error
	Update(ctx context.Context, collection, documentId string, changes map[string]any) (documentmodels.Document, error)
}

type ServerAPI struct {
	dbv1.UnimplementedDocumentServiceServer
	service DocumentService
}

func New(service DocumentService) *ServerAPI {
	return &ServerAPI{service: service}
}

func Register(server *grpc.Server, service DocumentService) {
	dbv1.RegisterDocumentServiceServer(server, New(service))
}

func (s *ServerAPI) CreateDocument(ctx context.Context, req *dbv1.CreateDocumentRequest) (*dbv1.Document, error) {
	return nil, nil
}

func (s *ServerAPI) DeleteDocument(ctx context.Context, req *dbv1.DeleteDocumentRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *ServerAPI) GetDocument(ctx context.Context, req *dbv1.GetDocumentRequest) (*dbv1.Document, error) {
	return nil, nil
}

func (s *ServerAPI) ListDocuments(ctx context.Context, req *dbv1.ListDocumentsRequest) (*dbv1.ListDocumentsResponse, error) {
	return nil, nil
}

func (s *ServerAPI) UpdateDocument(ctx context.Context, req *dbv1.UpdateDocumentRequest) (*dbv1.Document, error) {
	return nil, nil
}

func convert(src documentmodels.Document) (*dbv1.Document, error) {
	return nil, nil
}
