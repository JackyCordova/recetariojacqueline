// arranca el servicio de recipe, se registra en consul y envía heartbeats
// nos ayuda a descubrir los demás microservicios
package main

import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"recetariojacqueline.com/pkg/discovery/consul"
	discovery "recetariojacqueline.com/pkg/registry"
	"recetariojacqueline.com/recipe/internal/controller/recipe"
	metadatagateway "recetariojacqueline.com/recipe/internal/gateway/metadata/http"
	ratinggateway "recetariojacqueline.com/recipe/internal/gateway/rating/http"
	httphandler "recetariojacqueline.com/recipe/internal/handler/http"
)

const serviceName = "recipe" // nombre del servicio que se registra en consul y genera el id

// helper para leer variables de entorno con valor por defecto
func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() { // arrancamos el microservicio
	var port int
	flag.IntVar(&port, "port", 8083, "API Handler port") // definimos el puerto
	flag.Parse()
	log.Printf("Starting %s service on port: %d", serviceName, port)

	// dirección de consul desde la variable de entorno
	consulAddr := getenv("CONSUL_ADDR", "consul:8500")
	registry, err := consul.NewRegistry(consulAddr)
	if err != nil {
		panic(err)
	}

	ctx := context.Background() // contexto necesario

	// host del servicio desde la variable de entorno
	serviceHost := getenv("SERVICE_HOST", "recipe")

	// generamos el id único y lo registramos
	instanceID := discovery.GenerateInstanceID(serviceName)
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("%s:%d", serviceHost, port)); err != nil {
		// se registra en consul y permite que otros microservicios lo encuentren
		panic(err)
	}

	// rutina que reporta cada segundo al consul que está vivo
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName) // cuando se apaga, se elimina el registro del consul

	// llamada a otros microservicios
	metadataGateway := metadatagateway.New(registry)
	ratingGateway := ratinggateway.New(registry)

	// controller y handler HTTP
	ctrl := recipe.New(metadataGateway, ratingGateway) // respeta el orden del constructor
	h := httphandler.New(ctrl)

	// ruta
	http.Handle("/recipe", http.HandlerFunc(h.GetRecipeDetails))

	// iniciamos el servidor http y escucha peticiones
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
