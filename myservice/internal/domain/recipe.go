package domain

type Recipe struct {
	ID          int64
	Name        string
	Ingredients []RecipeIngredientLine
}

type RecipeIngredientLine struct {
	IngredientID int64
	Name         string
	Unit         string
	Quantity     int
}