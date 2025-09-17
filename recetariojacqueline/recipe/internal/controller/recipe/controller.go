// aqui es donde juntamos el metadata y rating para entregar la receta
package recipe //en este paquete vive la lógica del servicio de receta completas

import (
	"context"
	"errors"

	"recetariojacqueline.com/recipe/internal/gateway"
	model "recetariojacqueline.com/recipe/pkg/model"
)

// interfaces que necesita
type metadataRepo interface { //implementaciónes que llaman a get
	Get(ctx context.Context, id string) (map[string]any, error)
}
type ratingRepo interface { // implementaciones que tienen getaverage
	GetAverage(ctx context.Context, id string) (float64, int, error)
}

type Controller struct { // estructura principal del controller
	meta   metadataRepo //metadata
	rating ratingRepo   //info del rating
}

func (c *Controller) Get(context context.Context, id string) (any, any) {
	panic("unimplemented")
}

func New(meta metadataRepo, rating ratingRepo) *Controller { //constructor
	//recibe la metadata y la info del rating y nos devuelve el controller
	return &Controller{meta: meta, rating: rating}
}

func (c *Controller) GetRecipe(ctx context.Context, id string) (*model.Recipe, error) {
	//esto es lo que nos va a devolver una receta completa
	md, err := c.meta.Get(ctx, id)
	if err != nil {
		if errors.Is(err, gateway.ErrNotFound) { //si no existe la receta devolvemos error
			return nil, err
		}
		return nil, err //cualquier otro error que pueda ocurrir también marcamos lo que devuelve
	}

	// llamamos al servicio de rating para calcular el promedio y cantidad de ratings
	avg, count, err := c.rating.GetAverage(ctx, id)
	if err != nil && !errors.Is(err, gateway.ErrNotFound) {
		return nil, err
	}
	if err != nil && errors.Is(err, gateway.ErrNotFound) {
		avg, count = 0, 0 //si no hay info se pone 0
	}

	r := &model.Recipe{ //constuimos el objeto recipe
		//convertimos a variabes con un tipo específico - tipados
		ID:          toString(md["id"]),
		Title:       toString(md["title"]),
		Description: toString(md["description"]),
		Servings:    toInt(md["servings"]),
		Difficulty:  toString(md["difficulty"]),
		Ingredients: toStringSlice(md["ingredients"]),
		Utensils:    toStringSlice(md["utensils"]),
		Steps:       toStringSlice(md["steps"]),
		Average:     avg,
		Count:       count,
	}
	return r, nil //devolvemos la receta completa
}

// conversiones
func toString(v any) string {
	if s, ok := v.(string); ok {
		return s
	}
	return ""
}

func toInt(v any) int {
	switch x := v.(type) {
	case int:
		return x
	case float64:
		return int(x)
	default:
		return 0
	}
}

func toStringSlice(v any) []string {
	if ss, ok := v.([]string); ok {
		return ss
	}
	list, ok := v.([]any)
	if !ok {
		return nil
	}
	out := make([]string, 0, len(list))
	for _, it := range list {
		if s, ok := it.(string); ok {
			out = append(out, s)
		}
	}
	return out
}
