// estructuras de ratings
package model

type RecordID string
type RecordType string

const (
	RecordTypeRecipe = RecordType("recipe") //vamos a calificar recetas
)

type UserID string //usuario que da el rating

type RatingValue float64 //rating que le da

type Rating struct { //calificación del usuario
	RecordID   string      `json:"recordId"`   //receta calificada
	RecordType RecordType  `json:"recordType"` //tipo
	UserID     string      `json:"userId"`     // usuario
	Value      RatingValue `json:"value"`      //calificación
}

type Average struct { //promedio de las calificaciones
	RecordID string  `json:"recordId"` // receta calificada
	Avg      float64 `json:"avg"`      //promedio
	Count    int     `json:"count"`    //cuantos ratings tiene
}
