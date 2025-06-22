package collectionapi

import (
	"context"

	mdbv1 "github.com/10Narratives/distgo-db/pkg/proto/master/database/v1"
	dbv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/database/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type CollectionCoordinator interface {
	ListCollections(context.Context, *dbv1.ListCollectionsRequest) (*dbv1.ListCollectionsResponse, error)
	GetCollection(context.Context, *dbv1.GetCollectionRequest) (*dbv1.Collection, error)
	CreateCollection(context.Context, *dbv1.CreateCollectionRequest) (*dbv1.Collection, error)
	UpdateCollection(context.Context, *dbv1.UpdateCollectionRequest) (*dbv1.Collection, error)
	DeleteCollection(context.Context, *dbv1.DeleteCollectionRequest) (*emptypb.Empty, error)
}

type serverAPI struct {
	mdbv1.UnimplementedCollectionServiceServer
	coordinator CollectionCoordinator
}

func Register(server *grpc.Server, coordinator CollectionCoordinator) {
	mdbv1.RegisterCollectionServiceServer(server, &serverAPI{
		coordinator: coordinator,
	})
}

func (s *serverAPI) CreateCollection(ctx context.Context, req *mdbv1.CreateCollectionRequest) (*mdbv1.Collection, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.coordinator.CreateCollection(ctx, &dbv1.CreateCollectionRequest{
		Parent:       req.GetParent(),
		CollectionId: req.GetCollectionId(),
		Collection:   convertCollectionFromGRPC(req.GetCollection()),
	})

	return convertCollectionToGRPC(resp), err
}

func (s *serverAPI) DeleteCollection(ctx context.Context, req *mdbv1.DeleteCollectionRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	return s.coordinator.DeleteCollection(ctx, &dbv1.DeleteCollectionRequest{
		Name: req.GetName(),
	})
}

func (s *serverAPI) GetCollection(ctx context.Context, req *mdbv1.GetCollectionRequest) (*mdbv1.Collection, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.coordinator.GetCollection(ctx, &dbv1.GetCollectionRequest{
		Name: req.GetName(),
	})

	return convertCollectionToGRPC(resp), err
}

func (s *serverAPI) ListCollections(ctx context.Context, req *mdbv1.ListCollectionsRequest) (*mdbv1.ListCollectionsResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.coordinator.ListCollections(ctx, &dbv1.ListCollectionsRequest{
		Parent:    req.GetParent(),
		PageSize:  req.GetPageSize(),
		PageToken: req.GetPageToken(),
	})

	listed := make([]*mdbv1.Collection, 0, len(resp.Collections))
	for _, collection := range resp.Collections {
		listed = append(listed, convertCollectionToGRPC(collection))
	}

	return &mdbv1.ListCollectionsResponse{
		Collections:   listed,
		NextPageToken: resp.GetNextPageToken(),
	}, err
}

func (s *serverAPI) UpdateCollection(ctx context.Context, req *mdbv1.UpdateCollectionRequest) (*mdbv1.Collection, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	resp, err := s.coordinator.UpdateCollection(ctx, &dbv1.UpdateCollectionRequest{
		Collection: convertCollectionFromGRPC(req.GetCollection()),
		UpdateMask: req.GetUpdateMask(),
	})

	return convertCollectionToGRPC(resp), err
}

func convertCollectionFromGRPC(src *mdbv1.Collection) *dbv1.Collection {
	if src == nil {
		return nil
	}

	return &dbv1.Collection{
		Name:        src.GetName(),
		Description: src.GetDescription(),
		CreatedAt:   src.GetCreatedAt(),
		UpdatedAt:   src.GetUpdatedAt(),
	}
}

func convertCollectionToGRPC(src *dbv1.Collection) *mdbv1.Collection {
	if src == nil {
		return nil
	}

	return &mdbv1.Collection{
		Name:        src.GetName(),
		Description: src.GetDescription(),
		CreatedAt:   src.GetCreatedAt(),
		UpdatedAt:   src.GetUpdatedAt(),
	}
}
