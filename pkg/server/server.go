package server

import (
	"context"
	"fmt"
	"net"
	"net/http"

	grpcruntime "github.com/grpc-ecosystem/grpc-gateway/runtime"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/services"
	"github.com/sourcenetwork/acp_core/pkg/types"
)

type Server struct {
	server   *grpc.Server
	listener *net.TCPListener
	addr     string
}

func NewServer(address string) Server {
	server := grpc.NewServer()
	return Server{
		addr:   address,
		server: server,
	}
}

func (s *Server) Init(r runtime.RuntimeManager) error {
	engine := services.NewACPEngine(r)
	types.RegisterACPEngineServer(s.server, engine)

	addr, err := net.ResolveTCPAddr("tcp", s.addr)
	if err != nil {
		return err
	}
	s.listener, err = net.ListenTCP("tcp", addr)
	if err != nil {
		return err
	}

	return nil
}

func (s *Server) Run() {
	s.server.Serve(s.listener)
}

func NewGRPCGatewayServer(ctx context.Context, grpcAddress, gatewayAddress string) (*http.Server, error) {
	conn, err := grpc.DialContext(
		ctx,
		grpcAddress,
		grpc.WithBlock(),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		return nil, fmt.Errorf("dial to gRPC server for gateway: %w", err)
	}

	mux := grpcruntime.NewServeMux()
	err = types.RegisterACPEngineHandler(ctx, mux, conn)
	if err != nil {
		return nil, err
	}

	opts := []grpc.DialOption{
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	}
	err = types.RegisterACPEngineHandlerFromEndpoint(ctx, mux, grpcAddress, opts)
	if err != nil {
		return nil, err
	}

	return &http.Server{
		Addr:    gatewayAddress,
		Handler: mux,
	}, nil
}
