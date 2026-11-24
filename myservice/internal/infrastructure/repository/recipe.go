package repository

import "github.com/jackc/pgx/v5/pgxpool"

type RecipeRepository struct {
	pool *pgxpool.Pool
}

func NewRecipeRepository(pool *pgxpool.Pool) *RecipeRepository {
	return &RecipeRepository{pool: pool}
}