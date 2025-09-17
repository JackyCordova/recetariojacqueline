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

func New() *Repository { // crea nuevo repositorio vacio
	return &Repository{data: map[string]*model.Metadata{}}
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
