package documentapi

import (
	"context"

	mdbv1 "github.com/10Narratives/distgo-db/pkg/proto/master/database/v1"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type DocumentCoordinator interface {
	ListDocuments(context.Context, *dbv1.ListDocumentsRequest) (*dbv1.ListDocumentsResponse, error)
	GetDocument(context.Context, *dbv1.GetDocumentRequest) (*dbv1.Document, error)
	CreateDocument(context.Context, *dbv1.CreateDocumentRequest) (*dbv1.Document, error)
	UpdateDocument(context.Context, *dbv1.UpdateDocumentRequest) (*dbv1.Document, error)
	DeleteDocument(context.Context, *dbv1.DeleteDocumentRequest) (*emptypb.Empty, error)
}

type serverAPI struct {
	mdbv1.UnimplementedDocumentServiceServer
	coordinator DocumentCoordinator
}

func Register(server *grpc.Server, coordinator DocumentCoordinator) {
	mdbv1.RegisterDocumentServiceServer(server, &serverAPI{
		coordinator: coordinator,
	})
}

func (s *serverAPI) CreateDocument(ctx context.Context, req *mdbv1.CreateDocumentRequest) (*mdbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.coordinator.CreateDocument(ctx, &dbv1.CreateDocumentRequest{
		Parent:     req.GetParent(),
		DocumentId: req.GetDocumentId(),
		Document:   convertDocumentFromGRPC(req.GetDocument()),
	})

	return convertDocumentToGRPC(resp), err
}

func (s *serverAPI) DeleteDocument(ctx context.Context, req *mdbv1.DeleteDocumentRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return s.coordinator.DeleteDocument(ctx, &dbv1.DeleteDocumentRequest{
		Name: req.GetName(),
	})
}

func (s *serverAPI) GetDocument(ctx context.Context, req *mdbv1.GetDocumentRequest) (*mdbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.coordinator.GetDocument(ctx, &dbv1.GetDocumentRequest{
		Name: req.GetName(),
	})

	return convertDocumentToGRPC(resp), err
}

func (s *serverAPI) ListDocuments(ctx context.Context, req *mdbv1.ListDocumentsRequest) (*mdbv1.ListDocumentsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.coordinator.ListDocuments(ctx, &dbv1.ListDocumentsRequest{
		Parent:    req.GetParent(),
		PageSize:  req.GetPageSize(),
		PageToken: req.GetPageToken(),
	})

	listed := make([]*mdbv1.Document, 0, len(resp.Documents))
	for _, document := range resp.Documents {
		listed = append(listed, convertDocumentToGRPC(document))
	}

	return &mdbv1.ListDocumentsResponse{
		Documents:     listed,
		NextPageToken: resp.GetNextPageToken(),
	}, err
}

func (s *serverAPI) UpdateDocument(ctx context.Context, req *mdbv1.UpdateDocumentRequest) (*mdbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.coordinator.UpdateDocument(ctx, &dbv1.UpdateDocumentRequest{
		Document:   convertDocumentFromGRPC(req.GetDocument()),
		UpdateMask: req.GetUpdateMask(),
	})

	return convertDocumentToGRPC(resp), err
}

func convertDocumentFromGRPC(src *mdbv1.Document) *dbv1.Document {
	if src == nil {
		return nil
	}

	return &dbv1.Document{
		Name:      src.GetName(),
		Id:        src.GetId(),
		Value:     src.GetValue(),
		CreatedAt: src.GetCreatedAt(),
		UpdatedAt: src.GetUpdatedAt(),
	}
}

func convertDocumentToGRPC(src *dbv1.Document) *mdbv1.Document {
	if src == nil {
		return nil
	}

	return &mdbv1.Document{
		Name:      src.GetName(),
		Id:        src.GetId(),
		Value:     src.GetValue(),
		CreatedAt: src.GetCreatedAt(),
		UpdatedAt: src.GetUpdatedAt(),
	}
}
