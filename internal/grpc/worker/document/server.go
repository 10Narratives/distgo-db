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
	CreateDocument(ctx context.Context, collection, content string) (documentmodels.Document, error)
	Document(ctx context.Context, collection, documentID string) (documentmodels.Document, error)
	Documents(ctx context.Context, collection string) ([]documentmodels.Document, error)
	DeleteDocument(ctx context.Context, collection, documentID string) error
	UpdateDocument(ctx context.Context, collection, documentID string, changes string) (documentmodels.Document, error)

	CreateCollection(ctx context.Context, collectionID string) (documentmodels.Collection, error)
	Collections(ctx context.Context, size int, token string) ([]documentmodels.Collection, string, int)
}

type serverAPI struct {
	dbv1.UnimplementedDocumentServiceServer
	service DocumentService
}

var _ dbv1.DocumentServiceServer = &serverAPI{}

func Register(server *grpc.Server, service DocumentService) {
	dbv1.RegisterDocumentServiceServer(server, &serverAPI{service: service})
}

func (s *serverAPI) CreateDocument(ctx context.Context, req *dbv1.CreateDocumentRequest) (*dbv1.Document, error) {
	return nil, nil
}

func (s *serverAPI) DeleteDocument(ctx context.Context, req *dbv1.DeleteDocumentRequest) (*emptypb.Empty, error) {
	return nil, nil
}

func (s *serverAPI) GetDocument(ctx context.Context, req *dbv1.GetDocumentRequest) (*dbv1.Document, error) {
	return nil, nil
}

func (s *serverAPI) ListDocuments(ctx context.Context, req *dbv1.ListDocumentsRequest) (*dbv1.ListDocumentsResponse, error) {
	return nil, nil
}

func (s *serverAPI) UpdateDocument(ctx context.Context, req *dbv1.UpdateDocumentRequest) (*dbv1.Document, error) {
	return nil, nil
}

func (s *serverAPI) CreateCollection(context.Context, *dbv1.CreateCollectionRequest) (*dbv1.Collection, error) {
	panic("unimplemented")
}

func (s *serverAPI) ListCollections(context.Context, *dbv1.ListCollectionsRequest) (*dbv1.ListCollectionsResponse, error) {
	panic("unimplemented")
}

func (s *serverAPI) BeginTransaction(context.Context, *dbv1.BeginTransactionRequest) (*dbv1.BeginTransactionResponse, error) {
	panic("unimplemented")
}

func (s *serverAPI) CommitTransaction(context.Context, *dbv1.CommitTransactionRequest) (*dbv1.CommitTransactionResponse, error) {
	panic("unimplemented")
}

func (s *serverAPI) RollbackTransaction(context.Context, *dbv1.RollbackTransactionRequest) (*emptypb.Empty, error) {
	panic("unimplemented")
}
