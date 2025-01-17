package server

import (
	"net"

	"github.com/sourcenetwork/acp_core/pkg/runtime"
	"github.com/sourcenetwork/acp_core/pkg/services"
	"github.com/sourcenetwork/acp_core/pkg/types"
	"google.golang.org/grpc"
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
