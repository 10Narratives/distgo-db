package databasegrpc

import (
	"context"

	"github.com/10Narratives/distgo-db/internal/models"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DatabaseService interface {
	ListDocuments(ctx context.Context, collection string) ([]models.Document, error)
	GetDocument(ctx context.Context, collection string, documentID string) (models.Document, error)
	CreateDocument(ctx context.Context, collection string, data map[string]any) (models.Document, error)
	UpdateDocument(ctx context.Context, collection string, update map[string]any) (models.Document, error)
	DeleteDocument(ctx context.Context, collection, documentID string) (bool, error)
}

type serverAPI struct {
	dbv1.UnimplementedDatabaseServiceServer
	srv DatabaseService
}

func Register(gRPCServer *grpc.Server, dbSrv DatabaseService) {
	dbv1.RegisterDatabaseServiceServer(gRPCServer, &serverAPI{srv: dbSrv})
}

func (s *serverAPI) CreateDocument(context.Context, *dbv1.CreateDocumentRequest) (*dbv1.Document, error) {
	panic("implement me")
}

func (s *serverAPI) DeleteDocument(context.Context, *dbv1.DeleteDocumentRequest) (*emptypb.Empty, error) {
	panic("implement me")
}

func (s *serverAPI) GetDocument(context.Context, *dbv1.GetDocumentRequest) (*dbv1.Document, error) {
	panic("implement me")
}

func (s *serverAPI) ListDocuments(context.Context, *dbv1.ListDocumentsRequest) (*dbv1.ListDocumentsResponse, error) {
	panic("implement me")
}

func (s *serverAPI) UpdateDocument(context.Context, *dbv1.UpdateDocumentRequest) (*dbv1.Document, error) {
	panic("implement me")
}
