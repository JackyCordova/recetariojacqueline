package main

import (
	"log"
	"net"
	"os"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"recetariojacqueline.com/recipe/internal/controller/recipe"
	grpcserver "recetariojacqueline.com/recipe/internal/handler/grpc"
	"recetariojacqueline.com/src/gen"

	metadatagrpc "recetariojacqueline.com/recipe/internal/gateway/metadata/grpc"
	ratinggrpc "recetariojacqueline.com/recipe/internal/gateway/rating/grpc"
)

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	// Direcciones de otros microservicios (igual que BOOKS_ADDR/METADATA_ADDR en authors-grpc)
	metadataAddr := getenv("METADATA_ADDR", "metadata-grpc:9091")
	ratingAddr := getenv("RATING_ADDR", "rating-grpc:9092")

	metaGW := metadatagrpc.New(metadataAddr)
	rateGW := ratinggrpc.New(ratingAddr)
	ctrl := recipe.New(metaGW, rateGW)

	addr := getenv("ADDR_GRPC", ":9093")
	log.Printf("[recipe grpc] listening on %s", addr)

	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("listen: %v", err)
	}

	s := grpc.NewServer()
	h := grpcserver.New(ctrl)

	gen.RegisterRecipeServiceServer(s, h)
	reflection.Register(s)

	log.Fatal(s.Serve(lis))
}
