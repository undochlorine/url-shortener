package server

import (
	"errors"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
	"google.golang.org/grpc/reflection"
	"log"
	"net"
	"url-shortener/common/validate"

	grpcmiddleware "github.com/grpc-ecosystem/go-grpc-middleware"
	"google.golang.org/grpc/health"
	"google.golang.org/grpc/health/grpc_health_v1"
	proto "url-shortener/pb/shortener"
	"url-shortener/services/shortener/config"
	"url-shortener/services/shortener/internal/service"
)

type Server struct {
	proto.UnimplementedShortenerServer
	service  service.Interface
	server   *grpc.Server
	listener net.Listener
}

type GrpcServerBuilder struct {
	options                   []grpc.ServerOption
	enabledReflection         bool
	disableDefaultHealthCheck bool
	service                   service.Interface
}

func (s *Server) GetListener() net.Listener {
	return s.listener
}

// DialOption configures how we set up the connection
func (sb *GrpcServerBuilder) AddOption(o grpc.ServerOption) {
	sb.options = append(sb.options, o)
}

func (sb *GrpcServerBuilder) EnableReflection() {
	sb.enabledReflection = true
}

func (sb *GrpcServerBuilder) DisableDefaultHealthCheck() {
	sb.disableDefaultHealthCheck = false
}

func (sb *GrpcServerBuilder) SetServerParameters(serverParams keepalive.ServerParameters) {
	keepAlive := grpc.KeepaliveParams(serverParams)
	sb.AddOption(keepAlive)
}

func (sb *GrpcServerBuilder) SetStreamInterceptors(interceptors []grpc.StreamServerInterceptor) {
	chain := grpc.StreamInterceptor(grpcmiddleware.ChainStreamServer(interceptors...))
	sb.AddOption(chain)
}

func (sb *GrpcServerBuilder) SetUnaryInterceptors(interceptors []grpc.UnaryServerInterceptor) {
	chain := grpc.UnaryInterceptor(grpcmiddleware.ChainUnaryServer(interceptors...))
	sb.AddOption(chain)
}

// todo: set logger

func (sb *GrpcServerBuilder) SetService(s service.Interface) {
	sb.service = s
}

func (sb *GrpcServerBuilder) Build(cfg *config.Config) (*Server, error) {
	server := grpc.NewServer(sb.options...)
	if !sb.disableDefaultHealthCheck {
		grpc_health_v1.RegisterHealthServer(server, health.NewServer())
	}

	reflection.Register(server)
	if cfg.AppMode != config.APP_MODE_PRODUCTION {
		sb.EnableReflection()
	}

	if sb.service == nil {
		return nil, fmt.Errorf("insufficient args to set on builder")
	}

	s := &Server{
		service:  sb.service,
		server:   server,
		listener: nil,
	}

	proto.RegisterShortenerServer(server, s)
	return s, nil
}

// Start the GRPC server
func (s *Server) Start(addr string) error {
	var err error

	validate.EnableDefaultConfig()

	s.listener, err = net.Listen("tcp", addr)

	if err != nil {
		msg := fmt.Sprintf("Failed to listen: %v", err)
		return errors.New(msg)
	}

	go s.serv()

	fmt.Printf("Listening on %v\n", s.listener.Addr())
	return nil
}

func (s *Server) Cleanup() {
	fmt.Println("Stopping the server")
	s.server.GracefulStop()
	fmt.Println("Closing the listener")
	_ = s.listener.Close()
	fmt.Println("End of Program")
}

func (s *Server) serv() {
	if err := s.server.Serve(s.listener); err != nil {
		log.Fatalf("Failed to serve: %v", err)
	}
}
