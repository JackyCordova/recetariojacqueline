package rating

import (
	"context"

	"recetariojacqueline.com/rating/pkg/model"
)

type Repository interface {
	GetAverage(ctx context.Context, id model.RecordID, t model.RecordType) (float64, int, error)
	Put(ctx context.Context, id model.RecordID, t model.RecordType, user string, value float64) error
}

type Controller struct {
	repo Repository
}

// Constructor exportado
func New(repo Repository) *Controller {
	return &Controller{repo: repo}
}

// Wrappers usados por el handler gRPC
func (c *Controller) GetAverage(ctx context.Context, id model.RecordID, t model.RecordType) (float64, int, error) {
	return c.repo.GetAverage(ctx, id, t)
}

func (c *Controller) Put(ctx context.Context, id model.RecordID, t model.RecordType, userID string, value float64) error {
	return c.repo.Put(ctx, id, t, userID, value)
}
