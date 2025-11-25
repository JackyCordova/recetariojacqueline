package recipe

import (
	"context"

	"recetariojacqueline.com/src/gen"
)

// Interfaces que el controller necesita (las implementan los gateways gRPC):
type MetadataRepo interface {
	Get(ctx context.Context, id string) (*gen.Metadata, error)
}
type RatingRepo interface {
	GetAverage(ctx context.Context, id string) (float64, int, error)
}

type Controller struct {
	meta MetadataRepo
	rate RatingRepo
}

func New(meta MetadataRepo, rate RatingRepo) *Controller {
	return &Controller{meta: meta, rate: rate}
}

type Details struct {
	ID          string
	Title       string
	Description string
	Ingredients []string
	Utensils    []string
	Steps       []string
	Servings    int
	Difficulty  string
	Average     float64
}

// Agrega metadata + rating
func (c *Controller) GetRecipe(ctx context.Context, id string) (*Details, error) {
	md, err := c.meta.Get(ctx, id)
	if err != nil {
		return nil, err
	}

	avg, _, err := c.rate.GetAverage(ctx, id)
	if err != nil {
		return nil, err
	}

	return &Details{
		ID:          md.Id,
		Title:       md.Recipe.Title,
		Description: md.Recipe.Description,
		Ingredients: md.Recipe.Ingredients,
		Utensils:    md.Recipe.Utensils,
		Steps:       md.Recipe.Steps,
		Servings:    int(md.Recipe.Servings),
		Difficulty:  md.Recipe.Difficulty,
		Average:     avg,
	}, nil
}
