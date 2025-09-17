// recibe las peticiones http de los clientes y las traduce a llamadas al controller
package http

import (
	"encoding/json"
	"errors"
	"log"
	"net/http"
	"strconv"

	"recetariojacqueline.com/rating/internal/controller/rating"
	"recetariojacqueline.com/rating/internal/repository"
	model "recetariojacqueline.com/rating/pkg/model"
)

type Handler struct { //guardamos la referencia al controller de rating
	ctrl *rating.Controller
}

func New(ctrl *rating.Controller) *Handler { //recibimos el controller y devolvemos un nuevo handler
	return &Handler{ctrl}
}

func (h *Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	//funcion que busca manejar las peticiones http
	switch req.Method { //según el método que llegue
	case http.MethodGet: //get = consulta
		//lee los parámetros
		recordID := model.RecordID(req.URL.Query().Get("id"))
		recordType := model.RecordType(req.URL.Query().Get("type"))
		if recordType == "" {
			recordType = model.RecordTypeRecipe
		}
		avg, count, err := h.ctrl.GetAverage(req.Context(), recordID, recordType)
		if err != nil && errors.Is(err, repository.ErrNotFound) {
			w.WriteHeader(http.StatusNotFound)
			return
		}
		writeJSON(w, model.Average{RecordID: string(recordID), Avg: avg, Count: count}) //regresa el json con promedio y núm de ratings
	case http.MethodPut, http.MethodPost: //guardar rating
		_ = req.ParseForm() //procesamos datos
		//datos del formulario
		recordID := model.RecordID(req.FormValue("recordID"))
		recordType := model.RecordType(req.FormValue("recordType"))
		userID := model.UserID(req.FormValue("userID"))
		v, err := strconv.ParseFloat(req.FormValue("value"), 64)

		if err != nil {
			w.WriteHeader(http.StatusBadRequest)
			return
		}
		if recordType == "" {
			recordType = model.RecordTypeRecipe
		}
		if err := h.ctrl.PutRating(req.Context(), recordID, recordType, model.Rating{UserID: string(userID), Value: model.RatingValue(v)}); err != nil {
			log.Printf("Repository put error: %v\n", err)
			w.WriteHeader(http.StatusInternalServerError)
		}
	default: //si no es ninguno de los métodos especificados = error
		w.WriteHeader(http.StatusBadRequest)
	}
}

func writeJSON(w http.ResponseWriter, v any) { // de go pasa a json para escribirlo en http
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(v)
}
