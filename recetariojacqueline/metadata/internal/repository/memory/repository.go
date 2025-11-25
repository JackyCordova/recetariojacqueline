// almacena recetas en memoria - guardamos y buscamos recetas en el metadata
package memory

import (
	"context"
	"errors"
	"strings"
	"sync"

	"recetariojacqueline.com/metadata/internal/repository"
	model "recetariojacqueline.com/metadata/pkg"
)

type Repository struct { //almacen en memoria de las recetas
	sync.RWMutex //ayuda a evitar las race conditions
	data         map[string]*model.Metadata
}

func New() *Repository {
	r := &Repository{data: map[string]*model.Metadata{}}

	// Datos iniciales de ejemplo (pon los que quieras)
	r.data["r1"] = &model.Metadata{
		ID:          "r1",
		Title:       "Chilaquiles Verdes",
		Description: "Receta tradicional mexicana",
		Ingredients: []string{"Tortillas", "Salsa verde", "Crema"},
		Utensils:    []string{"Sartén"},
		Steps:       []string{"Cortar", "Freír", "Servir"},
		Servings:    2,
		Difficulty:  "Easy",
	}

	r.data["r2"] = &model.Metadata{
		ID:          "r2",
		Title:       "Pasta Alfredo",
		Description: "Cremosa y deliciosa",
		Ingredients: []string{"Pasta", "Crema", "Queso Parmesano"},
		Utensils:    []string{"Olla"},
		Steps:       []string{"Hervir", "Mezclar"},
		Servings:    4,
		Difficulty:  "Medium",
	}

	return r
}

func (r *Repository) Get(ctx context.Context, id string) (*model.Metadata, error) { // buscar una receta por su id
	r.RLock()                    //bloqueo de modo lectura
	defer r.RUnlock()            //desbloqueo
	if v, ok := r.data[id]; ok { //busca si existe el id
		return v, nil // devuelve la receta
	}
	return nil, repository.ErrNotFound //si no existe devuelve el error
}

func (r *Repository) Put(ctx context.Context, m *model.Metadata) error { // actualizar y guardar una receta
	r.Lock()         //bloqueo de modo lectura
	defer r.Unlock() //desbloqueo
	m.ID = strings.TrimSpace(m.ID)
	if m.ID == "" {
		return errors.New("id required") //si no hay id devuelve error
	}
	r.data[m.ID] = m //guardar o sobreescribir la receta
	return nil
}
