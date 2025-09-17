package memory // Paquete principal del archivo - registro de memoria

import (
	"context"
	"errors"
	"sync"
	"time"

	discovery "recetariojacqueline.com/pkg/registry"
)

type serviceName string
type instanceID string

type Registry struct { //guarda las direcciones de los servicios
	sync.RWMutex // candado que permite un solo acceso a la vez
	//se va a permitir un editor a la vez - evita race conditions
	serviceAddrs map[serviceName]map[instanceID]*serviceInstance // estructura del mapa
}

type serviceInstance struct { //guarda el host y la última vez que estuvo activo
	hostPort   string
	lastActive time.Time
}

func NewRegistry() *Registry { //crear un registro vacio
	return &Registry{serviceAddrs: map[serviceName]map[instanceID]*serviceInstance{}}
}

func (r *Registry) Register(ctx context.Context, serviceN string, instanceId string, hostPort string) error {
	//agregar o actualizar una instancia del servicio
	r.Lock()         // bloqueo
	defer r.Unlock() //desbloqueo
	sName := serviceName(serviceN)
	iID := instanceID(instanceId)
	if _, ok := r.serviceAddrs[sName]; !ok { //si aun no hay lo creamos
		r.serviceAddrs[sName] = make(map[instanceID]*serviceInstance)
	}
	r.serviceAddrs[sName][iID] = &serviceInstance{hostPort: hostPort, lastActive: time.Now()} //guardamos y actualizamos la instancia

	return nil
}

func (r *Registry) Deregister(ctx context.Context, instanceId string, serviceN string) error {
	//quitar una instancia
	r.Lock()
	defer r.Unlock()
	sName := serviceName(serviceN)
	iID := instanceID(instanceId)
	if _, ok := r.serviceAddrs[sName]; !ok { //si no existe no hay nada que borrar
		return nil
	}
	delete(r.serviceAddrs[sName], iID) //borramos la instancia
	return nil
}

func (r *Registry) ReportHealthyState(instanceId string, serviceN string) error {
	//marcar la instancia como viva
	r.Lock()
	defer r.Unlock()
	sName := serviceName(serviceN)
	iID := instanceID(instanceId)            //casteamos instance a ID
	if _, ok := r.serviceAddrs[sName]; !ok { //mandamos error si el servicio no está registrado
		return errors.New("Service is not registered yet")
	}
	if _, ok := r.serviceAddrs[sName][iID]; !ok { //mandamos error si la instancia del servicio no está registrada
		return errors.New("Service instance is not registered yet")
	}
	r.serviceAddrs[sName][iID].lastActive = time.Now()
	return nil
}

func (r *Registry) ServiceAddress(ctx context.Context, serviceN string) ([]string, error) {
	//obtenemos las direcciones activas
	r.RLock()
	defer r.RUnlock()
	sName := serviceName(serviceN)
	if len(r.serviceAddrs[sName]) == 0 { // si no hay instancias registradas devolvemos error de que no se encontró
		return nil, discovery.ErrNotFound
	}
	var res []string //registro de direcciones vivas
	for _, i := range r.serviceAddrs[sName] {
		if i.lastActive.Before(time.Now().Add(-5 * time.Second)) {
			//tiene que reportar vida cada 5 seg
			continue
		}
		res = append(res, i.hostPort) //agregamos dirección válida
	}
	return res, nil //devolver direcciones
}
