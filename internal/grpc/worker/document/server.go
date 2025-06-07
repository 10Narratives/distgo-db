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
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	doc, err := s.service.Create(ctx, req.GetParent(), req.Content.AsMap())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convert(doc)
}

func (s *ServerAPI) DeleteDocument(ctx context.Context, req *dbv1.DeleteDocumentRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.service.Delete(ctx, req.GetCollection(), req.GetDocumentId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *ServerAPI) GetDocument(ctx context.Context, req *dbv1.GetDocumentRequest) (*dbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	doc, err := s.service.Get(ctx, req.GetCollection(), req.GetDocumentId())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convert(doc)
}

func (s *ServerAPI) ListDocuments(ctx context.Context, req *dbv1.ListDocumentsRequest) (*dbv1.ListDocumentsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	docs, err := s.service.List(ctx, req.GetParent())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	listed := make([]*dbv1.Document, 0)
	for _, doc := range docs {
		converted, _ := convert(doc)
		listed = append(listed, converted)
	}

	return &dbv1.ListDocumentsResponse{Documents: listed}, nil
}

func (s *ServerAPI) UpdateDocument(ctx context.Context, req *dbv1.UpdateDocumentRequest) (*dbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	doc, err := s.service.Update(ctx, req.GetCollection(), req.GetDocumentId(), req.GetContent().AsMap())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convert(doc)
}

func convert(src documentmodels.Document) (*dbv1.Document, error) {
	content, err := structpb.NewStruct(src.Content)
	if err != nil {
		return &dbv1.Document{}, err
	}
	return &dbv1.Document{
		Name:      src.ID.String(),
		Content:   content,
		CreatedAt: timestamppb.New(src.CreatedAt),
		UpdatedAt: timestamppb.New(src.UpdatedAt),
	}, nil
}
