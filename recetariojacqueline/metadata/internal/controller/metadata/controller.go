// aquí recibimos requests del handler y le hablamos al repo para obtener datos o errores
package metadata

import (
	"context"
	"errors"

	model "recetariojacqueline.com/metadata/pkg" //struct del metadata
)

var ErrNotFound = errors.New("Not found") // error por si no se encuentra algo

type metadataRepository interface { //interfaz para el controlador
	Get(ctx context.Context, id string) (*model.Metadata, error) // buscamos receta por id
	Put(ctx context.Context, m *model.Metadata) error            // guardamos o actualizamos receta
}

type Controller struct { // implementación de la interfaz
	repo metadataRepository
}

func New(repo metadataRepository) *Controller {
	//recibimos un repositorio y devolvemos un controller que usaremos
	return &Controller{repo}
}

func (c *Controller) Get(ctx context.Context, id string) (*model.Metadata, error) {
	//método con el que vamos a obtener una receta
	return c.repo.Get(ctx, id)
}

func (c *Controller) Put(ctx context.Context, m *model.Metadata) error {
	//método para guardar o actualizar la receta
	return c.repo.Put(ctx, m)
}
