// recibe peticiones http y las traduce al control del metadata
package http

import (
	"encoding/json"
	"net/http"
	"strconv"
	"strings"

	"recetariojacqueline.com/metadata/internal/controller/metadata"
	model "recetariojacqueline.com/metadata/pkg"
)

type Handler struct { //conecta el cpontrolador con http - traducir requests en llamadas
	ctrl *metadata.Controller
}

func New(ctrl *metadata.Controller) *Handler { //creamos nuevo handler y lo pasamos al controlador
	return &Handler{ctrl}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	//punto de entrada para las peticiones http al servicio
	switch r.Method { //revisamos los métodos http que recibimos
	case http.MethodGet: //get - obtener metadata
		id := r.URL.Query().Get("id") //lee el parámetro id
		res, err := h.ctrl.Get(r.Context(), id)
		if err != nil {
			w.WriteHeader(http.StatusNotFound) //mensaje de error
			return
		}
		writeJSON(w, res) // si existe devuelve como json

	case http.MethodPut, http.MethodPost: //put o post - guardar metadat
		_ = r.ParseForm()     //procesar los datos recibidos a través del form
		m := &model.Metadata{ //creamos uno nuevo con los valores
			ID:          r.FormValue("id"),
			Title:       r.FormValue("title"),
			Description: r.FormValue("description"),
			Ingredients: splitList(r.FormValue("ingredients")),
			Utensils:    splitList(r.FormValue("utensils")),
			Steps:       splitList(r.FormValue("steps")),
			Difficulty:  r.FormValue("difficulty"),
		}
		//leemos y convertimos porciones a enteros
		if s := strings.TrimSpace(r.FormValue("servings")); s != "" {
			if v, err := strconv.Atoi(s); err == nil {
				m.Servings = v
			}
		}
		//llamada para guardar la receta
		if err := h.ctrl.Put(r.Context(), m); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		writeJSON(w, m) //devolvemos la receta guardada como json
	default: //si llega otro método devolvemos error
		w.WriteHeader(http.StatusBadRequest)
	}
}

func writeJSON(w http.ResponseWriter, v any) {
	//convertir cualquier valor a json, lo convertimos en la respuesta http
	w.Header().Set("Content-Type", "application/json") //respuesta tipo json
	_ = json.NewEncoder(w).Encode(v)
}

func splitList(s string) []string {
	//los elementos que traen comas, los separamos
	if strings.TrimSpace(s) == "" {
		return nil
	}
	parts := strings.Split(s, ",") // divide por comas
	for i := range parts {
		parts[i] = strings.TrimSpace(parts[i])
	} //quita espacios extra
	return parts //devolvemos la lista
}
