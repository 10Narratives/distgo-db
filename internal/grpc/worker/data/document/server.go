package documentgrpc

// import (
// 	"context"
// 	"errors"
// 	"strings"

// 	documentmodels "github.com/10Narratives/distgo-db/internal/models/worker/document"
// 	workerstore "github.com/10Narratives/distgo-db/internal/storages/worker"
// 	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
// 	"github.com/google/uuid"
// 	"google.golang.org/grpc"
// 	"google.golang.org/grpc/codes"
// 	"google.golang.org/grpc/status"
// 	"google.golang.org/protobuf/types/known/emptypb"
// 	"google.golang.org/protobuf/types/known/timestamppb"
// )

// //go:generate mockery --name DocumentService --output ./mocks/
// type DocumentService interface {
// 	CreateDocument(ctx context.Context, collectionID, content string) (documentmodels.Document, error)
// 	Document(ctx context.Context, collectionID string, documentID uuid.UUID) (documentmodels.Document, error)
// 	Documents(ctx context.Context, collectionID string, size int, token string) ([]documentmodels.Document, string, int, error)
// 	DeleteDocument(ctx context.Context, collectionID string, documentID uuid.UUID) error
// 	UpdateDocument(ctx context.Context, collectionID string, documentID uuid.UUID, changes string) (documentmodels.Document, error)

// 	CreateCollection(ctx context.Context, collectionID string) (documentmodels.Collection, error)
// 	Collection(ctx context.Context, collectionID string) (documentmodels.Collection, error)
// 	Collections(ctx context.Context, size int, token string) ([]documentmodels.Collection, string, int)
// }

// type serverAPI struct {
// 	dbv1.UnimplementedDocumentServiceServer
// 	service DocumentService
// }

// var _ dbv1.DocumentServiceServer = &serverAPI{}

// func Register(server *grpc.Server, service DocumentService) {
// 	dbv1.RegisterDocumentServiceServer(server, &serverAPI{service: service})
// }

// func (s *serverAPI) CreateDocument(ctx context.Context, req *dbv1.CreateDocumentRequest) (*dbv1.Document, error) {
// 	if err := req.Validate(); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid create document request: "+err.Error())
// 	}

// 	_, err := s.service.Collection(ctx, req.GetParent())
// 	if errors.Is(err, workerstore.ErrCollectionNotFound) {
// 		return nil, status.Error(codes.NotFound, "cannot find required collection: "+err.Error())
// 	} else if err != nil {
// 		return nil, status.Error(codes.Internal, "cannot create document: "+err.Error())
// 	}

// 	document, err := s.service.CreateDocument(ctx, req.GetParent(), req.Document.GetContent())
// 	if err != nil {
// 		return nil, status.Error(codes.Internal, "cannot create document: "+err.Error())
// 	}

// 	return convert(document), nil
// }

// func (s *serverAPI) DeleteDocument(ctx context.Context, req *dbv1.DeleteDocumentRequest) (*emptypb.Empty, error) {
// 	if err := req.Validate(); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid create document request: "+err.Error())
// 	}

// 	collectionID, documentID, err := splitName(req.GetName())
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "document id is invalid: "+err.Error())
// 	}

// 	err = s.service.DeleteDocument(ctx, collectionID, documentID)
// 	if errors.Is(err, workerstore.ErrCollectionNotFound) {
// 		return nil, status.Error(codes.NotFound, "cannot find required collection")
// 	} else if errors.Is(err, workerstore.ErrDocumentNotFound) {
// 		return nil, status.Error(codes.NotFound, "cannot find required document")
// 	} else if err != nil {
// 		return nil, status.Error(codes.Internal, "cannot get document: "+err.Error())
// 	}

// 	return &emptypb.Empty{}, nil
// }

// func (s *serverAPI) GetDocument(ctx context.Context, req *dbv1.GetDocumentRequest) (*dbv1.Document, error) {
// 	if err := req.Validate(); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid create document request: "+err.Error())
// 	}

// 	collectionID, documentID, err := splitName(req.GetName())
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "document id is invalid: "+err.Error())
// 	}

// 	document, err := s.service.Document(ctx, collectionID, documentID)
// 	if errors.Is(err, workerstore.ErrCollectionNotFound) {
// 		return nil, status.Error(codes.NotFound, "cannot find required collection")
// 	} else if errors.Is(err, workerstore.ErrDocumentNotFound) {
// 		return nil, status.Error(codes.NotFound, "cannot find required document")
// 	} else if err != nil {
// 		return nil, status.Error(codes.Internal, "cannot get document: "+err.Error())
// 	}

// 	return convert(document), nil
// }

// func (s *serverAPI) ListDocuments(ctx context.Context, req *dbv1.ListDocumentsRequest) (*dbv1.ListDocumentsResponse, error) {
// 	if err := req.Validate(); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid create document request: "+err.Error())
// 	}

// 	listed, token, tSize, err := s.service.Documents(ctx, req.GetParent(), int(req.GetPageSize()), req.GetPageToken())
// 	if errors.Is(err, workerstore.ErrCollectionNotFound) {
// 		return nil, status.Error(codes.NotFound, "cannot find required collection")
// 	} else if err != nil {
// 		return nil, status.Error(codes.Internal, "cannot list documents: "+err.Error())
// 	}

// 	converted := make([]*dbv1.Document, 0, len(listed))
// 	for _, document := range listed {
// 		converted = append(converted, convert(document))
// 	}

// 	return &dbv1.ListDocumentsResponse{
// 		Documents:     converted,
// 		NextPageToken: token,
// 		TotalSize:     int32(tSize),
// 	}, nil
// }

// func (s *serverAPI) UpdateDocument(ctx context.Context, req *dbv1.UpdateDocumentRequest) (*dbv1.Document, error) {
// 	if err := req.Validate(); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "invalid create document request: "+err.Error())
// 	}

// 	collectionID, documentID, err := splitName(req.GetDocument().GetName())
// 	if err != nil {
// 		return nil, status.Error(codes.InvalidArgument, "document id is invalid: "+err.Error())
// 	}

// 	_, err = s.service.UpdateDocument(ctx, collectionID, documentID, req.GetDocument().GetContent())
// 	if errors.Is(err, workerstore.ErrCollectionNotFound) {
// 		return nil, status.Error(codes.NotFound, "cannot find required collection")
// 	} else if errors.Is(err, workerstore.ErrDocumentNotFound) {
// 		return nil, status.Error(codes.NotFound, "cannot find required document")
// 	} else if err != nil {
// 		return nil, status.Error(codes.Internal, "cannot update document: "+err.Error())
// 	}

// 	return nil, nil
// }

// func (s *serverAPI) CreateCollection(ctx context.Context, req *dbv1.CreateCollectionRequest) (*dbv1.Collection, error) {
// 	if err := req.Validate(); err != nil {
// 		return nil, status.Error(codes.InvalidArgument, err.Error())
// 	}

// 	collection, err := s.service.CreateCollection(ctx, req.GetCollectionId())
// 	if errors.Is(err, workerstore.ErrCollectionNotFound) {
// 		return nil, status.Error(codes.NotFound, "cannot find required collection")
// 	} else if err != nil {
// 		return nil, status.Error(codes.Internal, "cannot create collection"+err.Error())
// 	}

// 	return &dbv1.Collection{
// 		Name: collection.Name,
// 	}, nil
// }

// func (s *serverAPI) ListCollections(context.Context, *dbv1.ListCollectionsRequest) (*dbv1.ListCollectionsResponse, error) {
// 	panic("unimplemented")
// }

// func (s *serverAPI) BeginTransaction(context.Context, *dbv1.BeginTransactionRequest) (*dbv1.BeginTransactionResponse, error) {
// 	panic("unimplemented")
// }

// func (s *serverAPI) CommitTransaction(context.Context, *dbv1.CommitTransactionRequest) (*dbv1.CommitTransactionResponse, error) {
// 	panic("unimplemented")
// }

// func (s *serverAPI) RollbackTransaction(context.Context, *dbv1.RollbackTransactionRequest) (*emptypb.Empty, error) {
// 	panic("unimplemented")
// }

// func splitName(name string) (string, uuid.UUID, error) {
// 	sliced := strings.Split(name, "/")
// 	collection, document := sliced[1], sliced[3]
// 	documentID, err := uuid.Parse(document)
// 	if err != nil {
// 		return "", uuid.UUID{}, err
// 	}
// 	return collection, documentID, nil
// }

// func convert(src documentmodels.Document) *dbv1.Document {
// 	return &dbv1.Document{
// 		Name:       src.ID.String(),
// 		Content:    src.Content,
// 		CreateTime: timestamppb.New(src.CreateTime),
// 		UpdateTime: timestamppb.New(src.UpdateTime),
// 	}
// }
