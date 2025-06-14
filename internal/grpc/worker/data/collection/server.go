package collectiongrpc

import (
	"context"

	collectionmodels "github.com/10Narratives/distgo-db/internal/models/worker/data/collection"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

//go:generate mockery --name CollectionService --output ./mocks/
type CollectionService interface {
	CreateCollection(ctx context.Context, parent, collectionID string) (collectionmodels.Collection, error)
	DeleteCollection(ctx context.Context, name string) error
	UpdateCollection(ctx context.Context, collection collectionmodels.Collection, paths []string) (collectionmodels.Collection, error)
	Collection(ctx context.Context, name string) (collectionmodels.Collection, error)
	Collections(ctx context.Context, parent string, size int32, token string) ([]collectionmodels.Collection, string, error)
}

type ServerAPI struct {
	dbv1.UnimplementedCollectionServiceServer
	service CollectionService
}

func New(service CollectionService) *ServerAPI {
	return &ServerAPI{
		service: service,
	}
}

func Register(server *grpc.Server, service CollectionService) {
	dbv1.RegisterCollectionServiceServer(server, New(service))
}

var _ dbv1.CollectionServiceServer = &ServerAPI{}

func (s *ServerAPI) CreateCollection(ctx context.Context, req *dbv1.CreateCollectionRequest) (*dbv1.Collection, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	parent := req.GetParent()
	collectionID := req.GetCollectionId()

	created, err := s.service.CreateCollection(ctx, parent, collectionID)
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertCollectionToGRPC(created), nil
}

func (s *ServerAPI) DeleteCollection(ctx context.Context, req *dbv1.DeleteCollectionRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	err := s.service.DeleteCollection(ctx, req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}

func (s *ServerAPI) GetCollection(ctx context.Context, req *dbv1.GetCollectionRequest) (*dbv1.Collection, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	collection, err := s.service.Collection(ctx, req.GetName())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertCollectionToGRPC(collection), nil
}

func (s *ServerAPI) ListCollections(ctx context.Context, req *dbv1.ListCollectionsRequest) (*dbv1.ListCollectionsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	collections, nextToken, err := s.service.Collections(ctx, req.GetParent(), req.GetPageSize(), req.GetPageToken())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	listed := make([]*dbv1.Collection, 0, len(collections))
	for _, c := range collections {
		listed = append(listed, convertCollectionToGRPC(c))
	}

	return &dbv1.ListCollectionsResponse{
		Collections:   listed,
		NextPageToken: nextToken,
	}, nil
}

func (s *ServerAPI) UpdateCollection(ctx context.Context, req *dbv1.UpdateCollectionRequest) (*dbv1.Collection, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	collection := collectionmodels.Collection{
		Name: req.GetCollection().GetName(),
	}

	updated, err := s.service.UpdateCollection(ctx, collection, req.GetUpdateMask().GetPaths())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return convertCollectionToGRPC(updated), nil
}

func convertCollectionToGRPC(src collectionmodels.Collection) *dbv1.Collection {
	return &dbv1.Collection{
		Name: src.Name,
	}
}
