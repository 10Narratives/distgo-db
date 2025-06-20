package workergrpc

import (
	"context"
	"fmt"
	"log/slog"
	"net"
	"strconv"

	clusterv1 "github.com/10Narratives/distgo-db/pkg/proto/master/cluster/v1"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
)

type App struct {
	log            *slog.Logger
	GRPCServer     *grpc.Server
	port           int
	clusterService clusterv1.ClusterServiceClient
	Master         *grpc.ClientConn
	WorkerID       string
}

func New(
	log *slog.Logger,
	port int,
	masterPort int,
	databaseName string,
) *App {
	loggingOpts := []logging.Option{
		logging.WithLogOnEvents(
			logging.PayloadReceived, logging.PayloadSent,
		),
	}

	recoveryOpts := []recovery.Option{
		recovery.WithRecoveryHandler(func(p interface{}) (err error) {
			log.Error("Recovered from panic", slog.Any("panic", p))
			return status.Errorf(codes.Internal, "internal error")
		}),
	}

	gRPCServer := grpc.NewServer(grpc.ChainUnaryInterceptor(
		recovery.UnaryServerInterceptor(recoveryOpts...),
		logging.UnaryServerInterceptor(InterceptorLogger(log), loggingOpts...),
	))

	app := &App{
		log:        log,
		GRPCServer: gRPCServer,
		port:       port,
	}
	app.mustRegister(databaseName, port, masterPort)

	return app
}

func InterceptorLogger(l *slog.Logger) logging.Logger {
	return logging.LoggerFunc(func(ctx context.Context, lvl logging.Level, msg string, fields ...any) {
		l.Log(ctx, slog.Level(lvl), msg, fields...)
	})
}

func (g *App) MustRun() {
	if err := g.Run(); err != nil {
		panic(err)
	}
}

func (g *App) Run() error {
	const op = "worker.grpcapp.Run"

	l, err := net.Listen("tcp", fmt.Sprintf(":%d", g.port))
	if err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	g.log.Info("grpc server started", slog.String("addr", l.Addr().String()))

	if err := g.GRPCServer.Serve(l); err != nil {
		return fmt.Errorf("%s: %w", op, err)
	}

	return nil
}

func (a *App) Stop() {
	const op = "worker.grpcapp.Stop"

	a.log.With(slog.String("op", op)).
		Info("stopping gRPC server", slog.Int("port", a.port))

	a.GRPCServer.GracefulStop()
	a.mustUnregister()
}

type Registrar func(server *grpc.Server, service interface{})

func (a *App) Register(reg Registrar, service interface{}) {
	reg(a.GRPCServer, service)
}

func (a *App) mustRegister(databaseName string, port, masterPort int) {
	target := strconv.Itoa(masterPort)
	conn, err := grpc.NewClient("localhost:"+target, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		panic("cannot connect to master: " + err.Error())
	}
	a.Master = conn

	target = strconv.Itoa(port)

	a.clusterService = clusterv1.NewClusterServiceClient(conn)

	resp, err := a.clusterService.Register(context.Background(), &clusterv1.RegisterRequest{
		DatabaseName: databaseName,
		Address:      "localhost:" + target,
	})
	if err != nil {
		panic("cannot connect to master: " + err.Error())
	}

	a.WorkerID = resp.WorkerId
}

func (a *App) mustUnregister() {
	_, err := a.clusterService.Unregister(context.Background(), &clusterv1.UnregisterRequest{
		WorkerId: a.WorkerID,
	})
	if err != nil {
		panic("cannot unregister from master: " + err.Error())
	}
}
