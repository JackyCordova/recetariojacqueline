package model

import "recetariojacqueline.com/src/gen"

func MetadataToProto(m *Metadata) *gen.Metadata {
	return &gen.Metadata{
		Id: m.ID,
		Recipe: &gen.RecipeStruct{
			Title:       m.Title,
			Description: m.Description,
			Ingredients: m.Ingredients,
			Utensils:    m.Utensils,
			Steps:       m.Steps,
			Servings:    int32(m.Servings),
			Difficulty:  m.Difficulty,
		},
	}
}

func MetadataFromProto(pm *gen.Metadata) *Metadata {
	if pm == nil || pm.Recipe == nil {
		return &Metadata{}
	}
	return &Metadata{
		ID:          pm.Id,
		Title:       pm.Recipe.Title,
		Description: pm.Recipe.Description,
		Ingredients: pm.Recipe.Ingredients,
		Utensils:    pm.Recipe.Utensils,
		Steps:       pm.Recipe.Steps,
		Servings:    int(pm.Recipe.Servings),
		Difficulty:  pm.Recipe.Difficulty,
	}
}
