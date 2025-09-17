// lo que le da lógica a mi servicio de ratings
package rating

import (
	"context"

	"recetariojacqueline.com/rating/pkg/model"
)

type ratingRepository interface { // interfaz del controlador
	Put(ctx context.Context, recID model.RecordID, userID model.UserID, val model.RatingValue) error // guardar el rating
	GetAverage(ctx context.Context, recID model.RecordID) (float64, int, error)                      //calcular el promedio y total de ratings
}

type Controller struct { // de que depende este controller
	repo ratingRepository
}

func New(repo ratingRepository) *Controller { //creamos un nuevo controller siempre y cuando cumpla con la interfaz
	return &Controller{repo}
}

func (c *Controller) PutRating(ctx context.Context, recID model.RecordID, _ model.RecordType, rating model.Rating) error {
	// guardamos el rating de la receta
	return c.repo.Put(ctx, recID, model.UserID(rating.UserID), rating.Value)
}

func (c *Controller) GetAverage(ctx context.Context, recID model.RecordID, _ model.RecordType) (float64, int, error) {
	// métodod para calcular el promedio de ratings
	return c.repo.GetAverage(ctx, recID)
}
