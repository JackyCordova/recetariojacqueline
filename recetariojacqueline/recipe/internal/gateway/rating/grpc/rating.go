package ratinggrpc

import (
	"context"
	"fmt"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"recetariojacqueline.com/src/gen"
)

// Gateway implementa recipe.RatingRepo
type Gateway struct {
	client gen.RatingServiceClient
}

// target: "rating-grpc:9092" en Kubernetes
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
		panic(fmt.Sprintf("rating grpc dial (%s): %v", target, err))
	}

	return &Gateway{
		client: gen.NewRatingServiceClient(conn),
	}
}

func (g *Gateway) GetAverage(ctx context.Context, id string) (float64, int, error) {
	res, err := g.client.GetAggregatedRating(ctx, &gen.GetAggregatedRatingRequest{
		RecordId:   id,
		RecordType: "recipe",
	})
	if err != nil {
		return 0, 0, err
	}
	return res.Avg, int(res.Count), nil
}
