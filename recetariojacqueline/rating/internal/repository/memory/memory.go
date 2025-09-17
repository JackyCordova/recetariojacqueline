// implementación de ratings
package memory

import (
	"context"
	"sync"

	"recetariojacqueline.com/rating/internal/repository"
	model "recetariojacqueline.com/rating/pkg/model"
)

type Repository struct {
	sync.RWMutex                                                       //candado para evitar race conditions
	userRatings  map[model.RecordID]map[model.UserID]model.RatingValue //mapa de la estructura
}

func New() *Repository { // inicializamos un repositorio vacio
	return &Repository{userRatings: map[model.RecordID]map[model.UserID]model.RatingValue{}}
}

func (r *Repository) Put(ctx context.Context, recID model.RecordID, userID model.UserID, val model.RatingValue) error {
	// guardar y actualizar el rating que un usuario le da a una receta
	r.Lock() //bloqueo
	defer r.Unlock()
	if _, ok := r.userRatings[recID]; !ok { // si no existe el id se crea
		r.userRatings[recID] = map[model.UserID]model.RatingValue{}
	}
	r.userRatings[recID][userID] = val // se guarda el rating del usuario
	return nil
}

func (r *Repository) GetAverage(ctx context.Context, recID model.RecordID) (float64, int, error) {
	// calculamos el promedio para una receta
	r.RLock() // bloqueo
	defer r.RUnlock()
	m := r.userRatings[recID] // jalamos todos los ratings que tiene esa receta
	if len(m) == 0 {
		return 0, 0, repository.ErrNotFound //si todavía no hay ratings devolvemos error
	}
	sum := 0.0
	for _, v := range m {
		sum += float64(v) //sumamos los ratings
	}
	return sum / float64(len(m)), len(m), nil //calculamos y devolvemos el promedio
}
