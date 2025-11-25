package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"

	"recetariojacqueline.com/rating/internal/controller/rating"
	grpcserver "recetariojacqueline.com/rating/internal/handler/grpc"
	"recetariojacqueline.com/rating/internal/repository/memory"
	"recetariojacqueline.com/src/gen"
)

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	addr := getenv("ADDR_GRPC", ":9092")

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	repo := memory.New()
	ctrl := rating.New(repo)
	h := grpcserver.New(ctrl)

	s := grpc.NewServer()
	gen.RegisterRatingServiceServer(s, h)

	log.Printf("[rating grpc] listening on %s", addr)
	log.Fatal(s.Serve(lis))
}
