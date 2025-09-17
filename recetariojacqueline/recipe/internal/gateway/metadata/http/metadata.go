package metadatagateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	discovery "recetariojacqueline.com/pkg/registry"
)

type Gateway struct { //guaramos la referencia al registro de servicios
	//preguntamos donde esta el servicio antes de que hagamos la petición
	registry discovery.Registry
}

// GetAverage implements recipe.ratingRepo.
func (g *Gateway) GetAverage(ctx context.Context, id string) (float64, int, error) {
	panic("unimplemented")
}

func New(reg discovery.Registry) *Gateway {
	//construimos gateway
	return &Gateway{registry: reg}
}

func (g *Gateway) Get(ctx context.Context, id string) (map[string]any, error) {
	// obtenemos la metadata de la receta por un id
	addrs, err := g.registry.ServiceAddress(ctx, "metadata")
	if err != nil || len(addrs) == 0 {
		return nil, err
	}
	//construimos la url del endpoint - nos aseguramos que se envíe bien el id
	endpoint := fmt.Sprintf("http://%s/metadata?id=%s", addrs[0], url.QueryEscape(id))
	// request get a metadata y cerramos la conexión antes de acabar
	resp, err := http.Get(endpoint)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 { //si no devuelve ok, mandamos el error
		return nil, fmt.Errorf("metadata %d", resp.StatusCode)
	}

	var out map[string]any
	//decodificamos el cuerpo de la respuesta
	if err := json.NewDecoder(resp.Body).Decode(&out); err != nil {
		return nil, err
	}
	return out, nil //devolvemos los datos de la receta
}
