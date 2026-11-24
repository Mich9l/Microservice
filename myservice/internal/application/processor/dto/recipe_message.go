package dto

type RecipeMessage struct {
	ID          string              `json:"id"`
	Name        string              `json:"name"`
	Description string              `json:"description"`
	Ingredients []IngredientMessage `json:"ingredients"`
}

type IngredientMessage struct {
	Name     string  `json:"name"`
	Unit     string  `json:"unit"`
	Quantity float64 `json:"quantity"`
}