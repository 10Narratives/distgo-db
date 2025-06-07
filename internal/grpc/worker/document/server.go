package documentgrpc

import (
	"context"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/structpb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name DocumentService --output ./mocks/
type DocumentService interface {
	Create(ctx context.Context, collection string, content map[string]any) (documentmodels.Document, error)
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
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	doc, err := s.service.Create(ctx, req.GetParent(), req.Content.AsMap())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	content, _ := structpb.NewStruct(doc.Content)
	return &dbv1.Document{
		Name:      doc.ID.String(),
		Content:   content,
		CreatedAt: timestamppb.New(doc.CreatedAt),
		UpdatedAt: timestamppb.New(doc.UpdatedAt),
	}, nil
}

func (s *ServerAPI) DeleteDocument(context.Context, *dbv1.DeleteDocumentRequest) (*emptypb.Empty, error) {
	panic("implement")

}

func (s *ServerAPI) GetDocument(context.Context, *dbv1.GetDocumentRequest) (*dbv1.Document, error) {
	panic("implement")

}

func (s *ServerAPI) ListDocuments(context.Context, *dbv1.ListDocumentsRequest) (*dbv1.ListDocumentsResponse, error) {
	panic("implement")

}

func (s *ServerAPI) UpdateDocument(context.Context, *dbv1.UpdateDocumentRequest) (*dbv1.Document, error) {
	panic("implement")
}
