// ayuda a que otros servicios consulten ratings
package ratinggateway

import (
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"

	discovery "recetariojacqueline.com/pkg/registry" //pregunta en qué direccipon está corriendo el servicio
)

type Gateway struct { //definimos la estructura
	registry discovery.Registry
}

// Get implements recipe.metadataRepo.
func (g *Gateway) Get(ctx context.Context, id string) (map[string]any, error) {
	panic("unimplemented")
}

func New(reg discovery.Registry) *Gateway { //construimos el gateway y lo devolvemos listo para usar
	return &Gateway{registry: reg}
}

func (g *Gateway) GetAverage(ctx context.Context, id string) (float64, int, error) { // consultamos el promedio
	// con registry preguntamos que está vivo
	addrs, err := g.registry.ServiceAddress(ctx, "rating")
	if err != nil || len(addrs) == 0 {
		return 0, 0, err
	}
	//construimos url a microservicio de rating
	endpoint := fmt.Sprintf("http://%s/rating?id=%s&type=recipe", addrs[0], url.QueryEscape(id))
	//request get al servicio, y despues de leer cerramos la conexxión
	resp, err := http.Get(endpoint)
	if err != nil {
		return 0, 0, err
	}
	defer resp.Body.Close()
	if resp.StatusCode != 200 { //si no devuelve ok, mostramos el error
		return 0, 0, fmt.Errorf("rating %d", resp.StatusCode)
	}
	var r struct {
		Avg   float64 `json:"avg"`
		Count int     `json:"count"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&r); err != nil {
		return 0, 0, err
	}

	return r.Avg, r.Count, nil // regresamos el promedio, num de ratings
}
