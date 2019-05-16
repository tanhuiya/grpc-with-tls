package main

import (
	"context"
	"log"
	"net"
	"runtime/debug"

	"github.com/grpc-ecosystem/go-grpc-middleware"

	"github.com/tanhuiya/grpc_with_tls/pkg/gtls"
	pb "github.com/tanhuiya/grpc_with_tls/proto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type SearchService struct{}

func (s *SearchService) Search(ctx context.Context, r *pb.SearchRequest) (*pb.SearchResponse, error) {
	if ctx.Err() == context.Canceled {
		return nil, status.Errorf(codes.Canceled, "searchService.Search canceled")
	}
	return &pb.SearchResponse{Response: r.GetRequest() + "HTTP Server"}, nil
}

const PORT = "9002"

func main() {
	// certFile := "../../cert/server/server.pem"
	// keyFile := "../../cert/server/server.key"
	// caFile := "../../cert/ca.pem"

	tlsServer := gtls.Server{
		CaFile:   "../../cert/ca.pem",
		CertFile: "../../cert/server/server.pem",
		KeyFile:  "../../cert/server/server.key",
	}

	c, err := tlsServer.GetCredentialByCA()
	if err != nil {
		log.Fatalf("GetTLSCredentialsByCA err: %v", err)
	}

	opts := []grpc.ServerOption{
		grpc.Creds(c),
		grpc_middleware.WithUnaryServerChain(
			RecoveryInterceptor,
			LoggingInterceptor,
		),
	}

	server := grpc.NewServer(opts...)
	pb.RegisterSearchServiceServer(server, &SearchService{})

	lis, err := net.Listen("tcp", ":"+PORT)
	if err != nil {
		log.Fatalf("net.Listen err: %v", err)
	}

	server.Serve(lis)
}

func LoggingInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	log.Printf("gRPC method: %s, %v", info.FullMethod, req)
	resp, err := handler(ctx, req)
	log.Printf("gRPC method: %s, %v", info.FullMethod, resp)
	return resp, err
}

func RecoveryInterceptor(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {
	defer func() {
		if e := recover(); e != nil {
			debug.PrintStack()
			err = status.Errorf(codes.Internal, "Panic err: %v", e)
		}
	}()

	return handler(ctx, req)
}
