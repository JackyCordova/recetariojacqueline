// punto de entrada del microservicio
// aquí es donde se levanta el servidor http en un puerto y registramos el servicio en consul
package main

// importar librerias y paquetes necesarios
import (
	"context"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"recetariojacqueline.com/metadata/internal/controller/metadata"
	httphandler "recetariojacqueline.com/metadata/internal/handler/http"
	"recetariojacqueline.com/metadata/internal/repository/memory"
	"recetariojacqueline.com/pkg/discovery/consul"
	discovery "recetariojacqueline.com/pkg/registry"
)

const serviceName = "metadata" // nombre del servicio

// helper para leer variables de entorno y si no existe devolver un valor por defecto
func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() { // función de punto de entrada
	var port int // puerto donde va a correr
	flag.IntVar(&port, "port", 8081, "API Handler port")
	flag.Parse()

	log.Printf("Starting %s service on port: %d", serviceName, port) //log que indica que el servicio arrancó en el puerto elegido

	// Dirección de Consul desde variable de entorno (default: consul:8500)
	consulAddr := getenv("CONSUL_ADDR", "consul:8500")
	registry, err := consul.NewRegistry(consulAddr)
	if err != nil {
		panic(err) // abortar si no se puede conectar
	}

	ctx := context.Background()
	instanceID := discovery.GenerateInstanceID(serviceName) // generamos un id único para la instancia

	// registra el servicio en consul con nombre y puerto
	serviceHost := getenv("SERVICE_HOST", "metadata")
	if err := registry.Register(ctx, instanceID, serviceName, fmt.Sprintf("%s:%d", serviceHost, port)); err != nil {
		panic(err)
	}

	// heartbeats a Consul cada segundo - nos ayuda a decir que el servicio esta vivo
	go func() {
		for {
			if err := registry.ReportHealthyState(instanceID, serviceName); err != nil {
				log.Println("Failed to report healthy state: " + err.Error())
			}
			time.Sleep(1 * time.Second)
		}
	}()
	defer registry.Deregister(ctx, instanceID, serviceName) //cuando el servicio termina se desregistra de consul

	// creamos todo y arrancamos el servidor http en el puerto
	repo := memory.New()
	ctrl := metadata.New(repo)
	h := httphandler.New(ctrl)
	http.Handle("/metadata", http.HandlerFunc(h.ServeHTTP))

	// servidor de http y se queda escuchando peticiones
	if err := http.ListenAndServe(fmt.Sprintf(":%d", port), nil); err != nil {
		panic(err) // si falla se detiene
	}
}
