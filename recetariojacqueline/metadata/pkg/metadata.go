package model

type Metadata struct {
	ID          string
	Title       string
	Description string
	Ingredients []string
	Utensils    []string
	Steps       []string
	Servings    int
	Difficulty  string
}
