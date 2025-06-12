package documentgrpc

import (
	"context"
	"fmt"

	"github.com/10Narratives/distgo-db/internal/lib/grpc/utils"
	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
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
	Create(ctx context.Context, collection, documentID, value string) (documentmodels.Document, error)
	Delete(ctx context.Context, collection, documentID string) error
	Document(ctx context.Context, collection, documentID string) (documentmodels.Document, error)
	Documents(ctx context.Context, collection string, pageSize int32, pageToken string) ([]documentmodels.Document, error)
	Update(ctx context.Context, collection, documentID, value string, paths []string) (documentmodels.Document, error)
}

//go:generate mockery --name CollectionService --output ./mocks/
type CollectionService interface {
	Collection(ctx context.Context, collection string) (collectionmodels.Collection, error)
}

type ServerAPI struct {
	dbv1.UnimplementedDocumentServiceServer

	documentSrv   DocumentService
	collectionSrv CollectionService
}

func New(documentSrv DocumentService, collectionSrv CollectionService) *ServerAPI {
	return &ServerAPI{
		documentSrv:   documentSrv,
		collectionSrv: collectionSrv,
	}
}

func Register(server *grpc.Server, documentSrv DocumentService, collectionSrv CollectionService) {
	dbv1.RegisterDocumentServiceServer(server, New(documentSrv, collectionSrv))
}

func (s *ServerAPI) CreateDocument(ctx context.Context, req *dbv1.CreateDocumentRequest) (*dbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	_, err := s.collectionSrv.Collection(ctx, req.GetParent())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	document, err := s.documentSrv.Create(ctx, req.GetParent(), req.GetDocumentId(), req.GetDocument().GetValue())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dbv1.Document{
		Name:      document.Name,
		Id:        document.ID,
		Value:     string(document.Value),
		CreatedAt: timestamppb.New(document.CreatedAt),
		UpdatedAt: timestamppb.New(document.UpdatedAt),
	}, nil
}

func (s *ServerAPI) DeleteDocument(ctx context.Context, req *dbv1.DeleteDocumentRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	n := utils.ParseName(req.GetName())

	_, err := s.collectionSrv.Collection(ctx, n.CollectionID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	err = s.documentSrv.Delete(ctx, n.CollectionID, n.DocumentID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *ServerAPI) GetDocument(ctx context.Context, req *dbv1.GetDocumentRequest) (*dbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	parsed := utils.ParseName(req.GetName())

	collection, err := s.collectionSrv.Collection(ctx, parsed.CollectionID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	document, err := s.documentSrv.Document(ctx, collection.Name, parsed.DocumentID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dbv1.Document{
		Name:      document.Name,
		Id:        document.ID,
		Value:     string(document.Value),
		CreatedAt: timestamppb.New(document.CreatedAt),
		UpdatedAt: timestamppb.New(document.UpdatedAt),
	}, nil
}

func (s *ServerAPI) ListDocuments(ctx context.Context, req *dbv1.ListDocumentsRequest) (*dbv1.ListDocumentsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	parsed := utils.ParseName(req.GetParent())
	fmt.Println(parsed)

	collection, err := s.collectionSrv.Collection(ctx, parsed.CollectionID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	docs, err := s.documentSrv.Documents(ctx, collection.Name, req.GetPageSize(), req.GetPageToken())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	responseDocs := make([]*dbv1.Document, len(docs))
	for i, doc := range docs {
		responseDocs[i] = &dbv1.Document{
			Name:      doc.Name,
			Id:        doc.ID,
			Value:     string(doc.Value),
			CreatedAt: timestamppb.New(doc.CreatedAt),
			UpdatedAt: timestamppb.New(doc.UpdatedAt),
		}
	}

	return &dbv1.ListDocumentsResponse{
		Documents:     responseDocs,
		NextPageToken: "",
	}, nil
}

func (s *ServerAPI) UpdateDocument(ctx context.Context, req *dbv1.UpdateDocumentRequest) (*dbv1.Document, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}
	parsed := utils.ParseName(req.GetDocument().GetName())

	collection, err := s.collectionSrv.Collection(ctx, parsed.CollectionID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	document, err := s.documentSrv.Update(ctx, collection.Name, parsed.DocumentID, req.Document.GetValue(), req.UpdateMask.Paths)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &dbv1.Document{
		Name:      document.Name,
		Id:        document.ID,
		Value:     string(document.Value),
		CreatedAt: timestamppb.New(document.CreatedAt),
		UpdatedAt: timestamppb.New(document.UpdatedAt),
	}, nil
}
