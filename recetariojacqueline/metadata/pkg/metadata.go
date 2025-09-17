// define la estructura de los datos de la receta
package model

type Metadata struct {
	//ayuda a que se convierta facilmente en json
	ID          string   `json:"id"`          //id de la receta
	Title       string   `json:"title"`       //nombre de la receta
	Description string   `json:"description"` //descripcion de la receta
	Ingredients []string `json:"ingredients"` //ingredientes para la receta
	Utensils    []string `json:"utensils"`    //utencilios
	Steps       []string `json:"steps"`       //pasos de la receta
	Servings    int      `json:"servings"`    //porciones para la receta
	Difficulty  string   `json:"difficulty"`  //nivel de dificultad
}
