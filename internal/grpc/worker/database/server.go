package databasegrpc

import (
	"context"

	protolib "github.com/10Narratives/distgo-db/internal/lib/proto"
	databasemodels "github.com/10Narratives/distgo-db/internal/models/worker/database"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type DatabaseService interface {
	CreateDocument(ctx context.Context, collection string, content map[string]any) (databasemodels.Document, error)
	Documents(ctx context.Context, collection string) ([]databasemodels.Document, error)
	Document(ctx context.Context, collection, documentID string) (databasemodels.Document, error)
	UpdateDocument(ctx context.Context, collection, documentID string) (databasemodels.Document, error)
	DeleteDocument(ctx context.Context, collection, documentID string) error
}

type serverAPI struct {
	dbv1.UnimplementedDatabaseServiceServer
	dbSrv DatabaseService
}

func Register(server *grpc.Server, service DatabaseService) {
	dbv1.RegisterDatabaseServiceServer(server, serverAPI{dbSrv: service})
}

func (s serverAPI) CreateDocument(ctx context.Context, req *dbv1.CreateDocumentRequest) (*dbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	doc, err := s.dbSrv.CreateDocument(ctx, req.GetCollection(), req.GetContent().AsMap())
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot create document")
	}

	return convertDocument(doc), nil
}

func (s serverAPI) DeleteDocument(ctx context.Context, req *dbv1.DeleteDocumentRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.dbSrv.DeleteDocument(ctx, req.GetCollection(), req.GetDocumentId())
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot delete document")
	}

	return &emptypb.Empty{}, nil
}

func (s serverAPI) GetDocument(ctx context.Context, req *dbv1.GetDocumentRequest) (*dbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	doc, err := s.dbSrv.Document(ctx, req.GetCollection(), req.GetDocumentId())
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot get document")
	}

	return convertDocument(doc), nil
}

func (s serverAPI) ListDocuments(ctx context.Context, req *dbv1.ListDocumentsRequest) (*dbv1.ListDocumentsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	docs, err := s.dbSrv.Documents(ctx, req.GetCollection())
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot list documents")
	}

	res := make([]*dbv1.Document, len(docs))
	for _, doc := range docs {
		res = append(res, convertDocument(doc))
	}

	return &dbv1.ListDocumentsResponse{Documents: res}, nil
}

func (s serverAPI) UpdateDocument(ctx context.Context, req *dbv1.UpdateDocumentRequest) (*dbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	doc, err := s.dbSrv.UpdateDocument(ctx, req.GetCollection(), req.GetDocumentId())
	if err != nil {
		return nil, status.Error(codes.Internal, "cannot update document")
	}

	return convertDocument(doc), nil
}

func convertDocument(src databasemodels.Document) *dbv1.Document {
	content, _ := protolib.MapToProtobufStruct(src.Content)
	return &dbv1.Document{
		Id:        src.ID.String(),
		Content:   content,
		CreatedAt: timestamppb.New(src.CreatedAt),
		UpdatedAt: timestamppb.New(src.UpdatedAt),
	}
}
