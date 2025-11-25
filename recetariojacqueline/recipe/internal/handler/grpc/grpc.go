package grpcserver

import (
	"context"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"recetariojacqueline.com/recipe/internal/controller/recipe"
	"recetariojacqueline.com/src/gen"
)

type Handler struct {
	gen.UnimplementedRecipeServiceServer
	ctrl *recipe.Controller
}

func New(ctrl *recipe.Controller) *Handler { return &Handler{ctrl: ctrl} }

func (h *Handler) GetRecipeDetails(ctx context.Context, req *gen.GetRecipeDetailsRequest) (*gen.GetRecipeDetailsResponse, error) {
	r, err := h.ctrl.GetRecipe(ctx, req.RecipeId)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "%v", err)
	}

	md := &gen.Metadata{
		Id: r.ID,
		Recipe: &gen.RecipeStruct{
			Title:       r.Title,
			Description: r.Description,
			Ingredients: r.Ingredients,
			Utensils:    r.Utensils,
			Steps:       r.Steps,
			Servings:    int32(r.Servings),
			Difficulty:  r.Difficulty,
		},
	}

	return &gen.GetRecipeDetailsResponse{
		RecipeDetails: &gen.RecipeDetails{
			Rating:   r.Average,
			Metadata: md,
		},
	}, nil
}
