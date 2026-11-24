package domain

type RecipeIngredient struct {
	ID           int64
	RecipeID     int64
	IngredientID int64
	Quantity     int
}
