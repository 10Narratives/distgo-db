package clustersrv

import (
	"context"
	"fmt"
	"strconv"

	clustergrpc "github.com/10Narratives/distgo-db/internal/grpc/worker/cluster"
	clusterv1 "github.com/10Narratives/distgo-db/pkg/proto/master/cluster/v1"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

type Service struct {
	masterConn *grpc.ClientConn
	master     clusterv1.ClusterServiceClient
	nodeID     string
}

var _ clustergrpc.ClusterService = &Service{}

func New(masterPort int) *Service {
	target := strconv.Itoa(masterPort)
	conn, err := grpc.NewClient("localhost:"+target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("cannot connect to master: " + err.Error())
	}

	return &Service{
		masterConn: conn,
		master:     clusterv1.NewClusterServiceClient(conn),
	}
}

func (s *Service) MustRegister(port int, nodeName string) {
	target := strconv.Itoa(port)
	resp, err := s.master.Register(context.Background(), &clusterv1.RegisterRequest{
		DatabaseName: nodeName,
		Address:      "localhost:" + target,
	})
	if err != nil {
		panic("cannot connect to master: " + err.Error())
	}
	s.nodeID = resp.WorkerId
	fmt.Println("nodeID", s.nodeID)
}

func (s *Service) Unregister() error {
	if s.masterConn.GetState() == connectivity.Shutdown {
		fmt.Println("Connection to master is already closed. Skipping unregistration.")
		return nil
	}

	_, err := s.master.Unregister(context.Background(), &clusterv1.UnregisterRequest{
		WorkerId: s.nodeID,
	})
	if err != nil {
		fmt.Println("Failed to unregister from master:", err)
		return err
	}

	fmt.Println("Successfully unregistered from master.")
	return nil
}

func (s *Service) ForgetConnection() {
	s.Unregister()
}
