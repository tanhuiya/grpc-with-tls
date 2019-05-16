package main

import (
	"context"
	"log"

	"github.com/tanhuiya/grpc_with_tls/pkg/gtls"
	"google.golang.org/grpc"

	pb "github.com/tanhuiya/grpc_with_tls/proto"
)

const PORT = "9002"

func main() {
	tlsClient := gtls.Client{
		ServerName: "grpc-with-tls",
		CertFile:   "../../cert/client/client.pem",
		KeyFile:    "../../cert/client/client.key",
		CaFile:     "../../cert/ca.pem",
	}
	c, err := tlsClient.GetCredentialByCA()
	if err != nil {
		log.Fatalf("tlsClient.GetTLSCredentials err: %v", err)
	}

	conn, err := grpc.Dial(":"+PORT, grpc.WithTransportCredentials(c))
	if err != nil {
		log.Fatalf("grpc.Dial err: %v", err)
	}
	defer conn.Close()

	client := pb.NewSearchServiceClient(conn)
	resp, err := client.Search(context.Background(), &pb.SearchRequest{
		Request: "gRPC",
	})
	if err != nil {
		log.Fatalf("client.Search err: %v", err)
	}

	log.Printf("resp: %s", resp.GetResponse())
}
