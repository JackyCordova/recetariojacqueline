// levanta el servidor, punto de entrada para ratings
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
	"recetariojacqueline.com/rating/internal/controller/rating"
	httphandler "recetariojacqueline.com/rating/internal/handler/http"
	"recetariojacqueline.com/rating/internal/repository/memory"
)

const serviceName = "rating" // nombre del servicio - cómo lo registro en consul

// helper para leer variables de entorno con valor por defecto
func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	var port int // número de puerto
	flag.IntVar(&port, "port", 8082, "API Handler port")
	flag.Parse()

	log.Printf("Starting %s service on port: %d", serviceName, port)

	// dirección de Consul desde la variable de entorno
	consulAddr := getenv("CONSUL_ADDR", "consul:8500")
	registry, err := consul.NewRegistry(consulAddr)
	if err != nil {
		panic(err) // abortar si no se puede conectar
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName) // id de la instancia

	// Host del servicio desde variable de entorno
	serviceHost := getenv("SERVICE_HOST", "rating")
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("%s:%d", serviceHost, port)); err != nil {
		panic(err) // registrar instancia o abortar si falla
	}

	// heartbeats a Consul cada segundo - dice que está vivo
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName) // eliminamos la instancia en cuanto termine el servicio

	// creamos todo - esto es para conectar todo
	repo := memory.New()
	ctrl := rating.New(repo)
	h := httphandler.New(ctrl)

	// asociamos la ruta al handler http
	http.Handle("/rating", http.HandlerFunc(h.ServeHTTP))

	// arrancamos el servidor - va a escuchar hasta que se detenga
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err)
	}
}
