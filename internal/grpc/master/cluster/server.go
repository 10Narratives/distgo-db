package clusterapi

import (
	"context"

	clustermodels "github.com/10Narratives/distgo-db/internal/models/master/cluster"
	clusterv1 "github.com/10Narratives/distgo-db/pkg/proto/master/cluster/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClusterService interface {
	Workers(ctx context.Context) []clustermodels.Worker
	Register(ctx context.Context, databaseName, address string) (string, error)
	Unregister(ctx context.Context, workerID string) error
}

type ServerAPI struct {
	clusterv1.UnimplementedClusterServiceServer
	clusterService ClusterService
}

func New(clusterService ClusterService) *ServerAPI {
	return &ServerAPI{
		clusterService: clusterService,
	}
}

func Register(server *grpc.Server, clusterService ClusterService) {
	clusterv1.RegisterClusterServiceServer(server, &ServerAPI{clusterService: clusterService})
}

func (s *ServerAPI) ListWorkers(ctx context.Context, req *clusterv1.ListWorkersRequest) (*clusterv1.ListWorkersResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	workers := s.clusterService.Workers(ctx)

	listed := make([]*clusterv1.Worker, 0, len(workers))
	for _, worker := range workers {
		listed = append(listed, &clusterv1.Worker{
			WorkerId:     worker.ID.String(),
			Address:      worker.Address,
			DatabaseName: worker.Database,
		})
	}

	return &clusterv1.ListWorkersResponse{
		Workers: listed,
	}, nil
}

func (s *ServerAPI) Register(ctx context.Context, req *clusterv1.RegisterRequest) (*clusterv1.RegisterResponse, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	workerID, err := s.clusterService.Register(ctx, req.GetDatabaseName(), req.GetAddress())
	if err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &clusterv1.RegisterResponse{
		WorkerId: workerID,
	}, nil
}

func (s *ServerAPI) Unregister(ctx context.Context, req *clusterv1.UnregisterRequest) (*emptypb.Empty, error) {
	if err := req.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	if err := s.clusterService.Unregister(ctx, req.GetWorkerId()); err != nil {
		return nil, status.Error(codes.Internal, err.Error())
	}

	return &emptypb.Empty{}, nil
}
