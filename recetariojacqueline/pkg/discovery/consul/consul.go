// implementación de registry
// ayuda a que los servicios puedan descubrirse entre sí
package consul

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"strings"

	consul "github.com/hashicorp/consul/api"
	discovery "recetariojacqueline.com/pkg/registry"
)

type Registry struct {
	client *consul.Client //cliente que habla con el agente Consul
}

func NewRegistry(addr string) (*Registry, error) { // creamos el registry con Consul
	cfg := consul.DefaultConfig()
	if addr != "" {
		cfg.Address = addr
	}
	client, err := consul.NewClient(cfg) //construir el cliente
	if err != nil {                      //si falla hay error
		return nil, err
	}
	return &Registry{client: client}, nil //devolver el registry con el cliente
}

func (r *Registry) Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error { // Definición de función
	parts := strings.Split(hostPort, ":") //espera el puerto
	if len(parts) != 2 {                  //validación del formato
		return errors.New("Hostport must be in a form of <host>:<port>, example: localhost:8081")
	}
	port, err := strconv.Atoi(parts[1]) //tomamos la parte numérica del puerto y se convierte de string a entero
	if err != nil {                     //validar si la conversion esta bien
		return err
	}
	return r.client.Agent().ServiceRegister(&consul.AgentServiceRegistration{
		//registro del servicio en el cliente
		Address: parts[0],
		ID:      instanceID,
		Name:    serviceName,
		Port:    port,
		Check: &consul.AgentServiceCheck{ //health check con time to live
			CheckID:                        instanceID,
			TTL:                            "5s",  //se reporta salud cada 5 seg
			DeregisterCriticalServiceAfter: "10s", //después de 10 seg se borra
		},
	})
}

func (r *Registry) Deregister(ctx context.Context, instanceID string, _ string) error {
	//borra la instancia por ID
	return r.client.Agent().ServiceDeregister(instanceID)
}

func (r *Registry) ServiceAddress(ctx context.Context, serviceID string) ([]string, error) {
	//devolver direcciones de las instancias sanas
	health := r.client.Health()                                 //consulta el estado de los servicios
	entries, _, err := health.Service(serviceID, "", true, nil) //consulta todas las instancias registradas
	if err != nil {                                             //error de consulta
		return nil, err
	} else if len(entries) == 0 { //error si no encontramos instancias sanas
		return nil, discovery.ErrNotFound
	}
	var res []string
	for _, e := range entries { //construir el hostport de cada instancia
		res = append(res, fmt.Sprintf("%s:%d", e.Service.Address, e.Service.Port))
	}
	return res, nil //devuelve las instancias vivas
}

func (r *Registry) ReportHealthyState(instanceID string, _ string) error {
	//avisa al consul que la instancia sigue viva - con el heartbeat
	return r.client.Agent().PassTTL(instanceID, "") //recibe el heartbeat, cada 5 segundos y llama al cliente consul
}
