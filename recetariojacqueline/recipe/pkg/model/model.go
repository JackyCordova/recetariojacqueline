// estructura de la receta
package model

type Recipe struct {
	ID          string   `json:"id"`          // id
	Title       string   `json:"title"`       // nombre de la receta
	Description string   `json:"description"` // descripción
	Ingredients []string `json:"ingredients"` //ingredientes a usar
	Utensils    []string `json:"utensils"`    //utensilios a usar
	Steps       []string `json:"steps"`       //pasos a seguir
	Servings    int      `json:"servings"`    // porciones o cantidades
	Difficulty  string   `json:"difficulty"`  //nivel de dificultad
	Average     float64  `json:"average"`     // promedio de ratings de la receta
	Count       int      `json:"count"`       // número de ratings que tiene
}
