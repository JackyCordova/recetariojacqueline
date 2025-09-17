// handler http para recetas
package http

import (
	"encoding/json"
	"net/http"

	"recetariojacqueline.com/recipe/internal/controller/recipe"
)

type Handler struct { // cada peticion http se traduce en una llamada al controller
	ctrl *recipe.Controller
}

func New(ctrl *recipe.Controller) *Handler { //constructor del handler
	return &Handler{ctrl}
}

func (h *Handler) GetRecipeDetails(w http.ResponseWriter, req *http.Request) {
	//métodos para manejar la ruta get de recipe
	id := req.URL.Query().Get("id") // extrae el id de la query string
	//llama al controler para obtener la receta con ese id
	res, err := h.ctrl.GetRecipe(req.Context(), id)
	if err != nil { //error si no encontramos la receta
		w.WriteHeader(http.StatusNotFound)
		return
	}
	//configuramos el header y convertimos a json
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(res)
}
