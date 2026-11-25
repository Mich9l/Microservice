package repository

import (
	"context"
	"database/sql"
	"fmt"
)

type RecipeRepository struct {
	db *sql.DB
}

func NewRecipeRepository(db *sql.DB) *RecipeRepository {
	return &RecipeRepository{db: db}
}

func (r *RecipeRepository) CreateRecipe(ctx context.Context, name string) (int64, error) {
	const query = `
		INSERT INTO recipes (name)
		VALUES ($1)
		RETURNING id
	`

	var id int64
	if err := r.db.QueryRowContext(ctx, query, name).Scan(&id); err != nil {
		return 0, fmt.Errorf("create recipe: %w", err)
	}

	return id, nil
}