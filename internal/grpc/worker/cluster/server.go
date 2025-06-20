package clustergrpc

import (
	"context"

	wclusterv1 "github.com/10Narratives/distgo-db/pkg/proto/worker/cluster/v1"
	"google.golang.org/grpc"
	"google.golang.org/protobuf/types/known/emptypb"
)

type ClusterService interface {
	ForgetConnection()
}

type ServerAPI struct {
	wclusterv1.UnimplementedClusterServiceServer
	clusterService ClusterService
}

func New(clusterService ClusterService) *ServerAPI {
	return &ServerAPI{
		clusterService: clusterService,
	}
}

func Register(server *grpc.Server, clusterService ClusterService) {
	wclusterv1.RegisterClusterServiceServer(server, &ServerAPI{clusterService: clusterService})
}

func (s *ServerAPI) ForgetConnection(context.Context, *wclusterv1.ForgetConnectionRequest) (*emptypb.Empty, error) {
	s.clusterService.ForgetConnection()
	return &emptypb.Empty{}, nil
}
