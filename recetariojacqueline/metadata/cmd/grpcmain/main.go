package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"recetariojacqueline.com/metadata/internal/controller/metadata"
	grpcserver "recetariojacqueline.com/metadata/internal/handler/grpc"
	"recetariojacqueline.com/metadata/internal/repository/memory"
	"recetariojacqueline.com/src/gen"
)

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	addr := getenv("ADDR_GRPC", ":9091")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	repo := memory.New()
	ctrl := metadata.New(repo)
	h := grpcserver.New(ctrl)

	s := grpc.NewServer()
	gen.RegisterMetadataServiceServer(s, h)

	log.Printf("[metadata grpc] listening on %s", addr)
	log.Fatal(s.Serve(lis))
}
