// Descubrir servicios dentro de la arquitectura de microservicios
package discovery // Paquete principal del archivo

import (
	"context"
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Registry interface { //ayuda a que un registro encuentre a otro
	//registrar un servicio con su dirección
	Register(ctx context.Context, instanceID string, serviceName string, hostPort string) error
	//si no está activo el servicio lo quitamos
	Deregister(ctx context.Context, instanceID string, serviceName string) error
	//pedir la dirección del servicio
	ServiceAddress(ctx context.Context, serviceID string) ([]string, error)
	// reporte de salud de que una instancia sigue viva - esto lo sabemos por el heartbeat
	ReportHealthyState(instanceID string, serviceName string) error
}

// error que devolvemos si no se encuentran direcciones de un servicio
var ErrNotFound = errors.New("No service addresses found")

func GenerateInstanceID(serviceName string) string { // crear un id único para cada instancia de servicio
	return fmt.Sprintf("%s-%d", serviceName, rand.New(rand.NewSource(time.Now().UnixNano())).Int())
}
