package server

import (
	"context"
	"net"

	"go.uber.org/zap"
	"google.golang.org/grpc"

	"myservice/internal/infrastructure/repository"
)

type GRPCServer struct {
	addr string
	srv  *grpc.Server
	log  *zap.Logger
}

func New(addr string, _ *repository.RecipeRepository, log *zap.Logger) (*GRPCServer, error) {
	return &GRPCServer{
		addr: addr,
		srv:  grpc.NewServer(),
		log:  log,
	}, nil
}

func (s *GRPCServer) Start() error {
	ln, err := net.Listen("tcp", s.addr)
	if err != nil {
		return err
	}
	return s.srv.Serve(ln)
}

func (s *GRPCServer) Shutdown(ctx context.Context) error {
	stopped := make(chan struct{})
	go func() {
		s.srv.GracefulStop()
		close(stopped)
	}()
	select {
	case <-ctx.Done():
		s.srv.Stop()
		return ctx.Err()
	case <-stopped:
		return nil
	}
}