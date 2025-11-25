package metadatagrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"recetariojacqueline.com/src/gen"
)

// Gateway implementa recipe.MetadataRepo
type Gateway struct {
	client gen.MetadataServiceClient
}

// target: "metadata-grpc:9091" en Kubernetes
func New(target string) *Gateway {
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	conn, err := grpc.DialContext(
		ctx,
		target,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
		grpc.WithBlock(),
	)
	if err != nil {
		panic(fmt.Sprintf("metadata grpc dial (%s): %v", target, err))
	}

	return &Gateway{
		client: gen.NewMetadataServiceClient(conn),
	}
}

func (g *Gateway) Get(ctx context.Context, id string) (*gen.Metadata, error) {
	res, err := g.client.GetMetadata(ctx, &gen.GetMetadataRequest{RecipeId: id})
	if err != nil {
		return nil, err
	}
	return res.Metadata, nil
}
