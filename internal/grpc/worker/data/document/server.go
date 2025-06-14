package documentgrpc

import (
	"context"
	"encoding/json"

	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/document"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"google.golang.org/protobuf/types/known/timestamppb"
)

//go:generate mockery --name DocumentService --output ./mocks/
type DocumentService interface {
	CreateDocument(ctx context.Context, parent string, documentID string, value string) (documentmodels.Document, error)
	DeleteDocument(ctx context.Context, name string) error
	UpdateDocument(ctx context.Context, document documentmodels.Document, paths []string) (documentmodels.Document, error)
	Document(ctx context.Context, name string) (documentmodels.Document, error)
	Documents(ctx context.Context, parent string, size int32, token string) ([]documentmodels.Document, string, error)
}

type ServerAPI struct {
	dbv1.UnimplementedDocumentServiceServer
	service DocumentService
}

func New(service DocumentService) *ServerAPI {
	return &ServerAPI{
		service: service,
	}
}

func Register(server *grpc.Server, service DocumentService) {
	dbv1.RegisterDocumentServiceServer(server, New(service))
}

var _ dbv1.DocumentServiceServer = &ServerAPI{}

func (s *ServerAPI) CreateDocument(ctx context.Context, req *dbv1.CreateDocumentRequest) (*dbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	parent := req.GetParent()
	documentID := req.GetDocumentId()
	Value := req.GetDocument().GetValue()

	doc, err := s.service.CreateDocument(ctx, parent, documentID, Value)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertDocumentToGRPC(doc), nil
}

func (s *ServerAPI) DeleteDocument(ctx context.Context, req *dbv1.DeleteDocumentRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.service.DeleteDocument(ctx, req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *ServerAPI) GetDocument(ctx context.Context, req *dbv1.GetDocumentRequest) (*dbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	doc, err := s.service.Document(ctx, req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertDocumentToGRPC(doc), nil
}

func (s *ServerAPI) ListDocuments(ctx context.Context, req *dbv1.ListDocumentsRequest) (*dbv1.ListDocumentsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	docs, nextToken, err := s.service.Documents(ctx, req.GetParent(), req.GetPageSize(), req.GetPageToken())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	listed := make([]*dbv1.Document, 0, len(docs))
	for _, d := range docs {
		listed = append(listed, convertDocumentToGRPC(d))
	}

	return &dbv1.ListDocumentsResponse{
		Documents:     listed,
		NextPageToken: nextToken,
	}, nil
}

func (s *ServerAPI) UpdateDocument(ctx context.Context, req *dbv1.UpdateDocumentRequest) (*dbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	document := documentmodels.Document{
		Name:  req.GetDocument().GetName(),
		Value: json.RawMessage(req.GetDocument().GetValue()),
	}

	updated, err := s.service.UpdateDocument(ctx, document, req.GetUpdateMask().GetPaths())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertDocumentToGRPC(updated), nil
}

func convertDocumentToGRPC(src documentmodels.Document) *dbv1.Document {
	return &dbv1.Document{
		Name:      src.Name,
		CreatedAt: timestamppb.New(src.CreatedAt),
		UpdatedAt: timestamppb.New(src.UpdatedAt),
		Value:     string(src.Value),
	}
}
