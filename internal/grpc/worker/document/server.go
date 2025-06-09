package documentgrpc

import (
	"context"
	"errors"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	workerstore "github.com/10Narratives/distgo-db/internal/storages/worker"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name DocumentService --output ./mocks/
type DocumentService interface {
	CreateDocument(ctx context.Context, collectionID, content string) (documentmodels.Document, error)
	Document(ctx context.Context, collectionID, documentID string) (documentmodels.Document, error)
	Documents(ctx context.Context, collectionID string, size int, token string) ([]documentmodels.Document, string, int, error)
	DeleteDocument(ctx context.Context, collectionID, documentID string) error
	UpdateDocument(ctx context.Context, collectionID, documentID string, changes string) (documentmodels.Document, error)

	CreateCollection(ctx context.Context, collectionID string) (documentmodels.Collection, error)
	Collection(ctx context.Context, collectionID string) (documentmodels.Collection, error)
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
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.Internal, "invalid create document request: "+err.Error())
	}

	_, err := s.service.Collection(ctx, req.GetParent())
	if errors.Is(err, workerstore.ErrCollectionNotFound) {
		return nil, status.Error(codes.NotFound, "cannot found required collection: "+err.Error())
	} else if err != nil {
		return nil, status.Error(codes.Internal, "cannot create document: "+err.Error())
	}

	document, err := s.service.CreateDocument(ctx, req.GetParent(), req.Document.GetContent())
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot create document: "+err.Error())
	}

	return &dbv1.Document{
		Name:       document.ID.String(),
		Content:    document.Content,
		CreateTime: timestamppb.New(document.CreateTime),
		UpdateTime: timestamppb.New(document.UpdateTime),
	}, nil
}

func (s *serverAPI) DeleteDocument(ctx context.Context, req *dbv1.DeleteDocumentRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.Internal, "invalid create document request: "+err.Error())
	}

	return nil, nil
}

func (s *serverAPI) GetDocument(ctx context.Context, req *dbv1.GetDocumentRequest) (*dbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.Internal, "invalid create document request: "+err.Error())
	}

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
