//go:generate protoc -I ../../api/ BannerRotationService.proto --go_out=. --go-grpc_out=.

package server

import (
	"context"
	"errors"
	"github.com/and67o/otus_project/internal/app"
	"github.com/and67o/otus_project/internal/configuration"
	"google.golang.org/grpc"
	"net"
)

const network = "tcp"

type Server struct {
	app       *app.App
	uBRServer UnimplementedBannerRotationServer
	addr      string
	server    *grpc.Server
}

type GRPC interface {
	Stop() error
	Start(ctx context.Context) error

	GRPCRoutes
}

func New(app *app.App, config configuration.GRPCConf) GRPC {
	return &Server{
		app:       app,
		uBRServer: UnimplementedBannerRotationServer{},
		addr:      net.JoinHostPort(config.Host, config.Port),
		server:    nil,
	}
}

func (s *Server) Stop() error {
	if s.server == nil {
		return errors.New("grpc server is nil")
	}

	s.server.GracefulStop()

	return nil
}

func (s *Server) Start(ctx context.Context) error {
	l, err := net.Listen(network, s.addr)
	if err != nil {
		return err
	}

	serverGRPC := grpc.NewServer()
	RegisterBannerRotationServer(serverGRPC, s.uBRServer)
	s.server = serverGRPC

	go func() {
		if err := s.server.Serve(l); err != nil {
			s.app.Logger.Error("error running server")
		}
	}()

	//err = s.server.Serve(l)
	//if err != nil {
	//	return err
	//}

	//<-ctx.Done()
	return nil
}
