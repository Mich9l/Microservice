package repository

import (
	"context"
	"database/sql"
	"fmt"

	"myservice/internal/domain"
)

type queryRower interface {
	QueryRowContext(ctx context.Context, query string, args ...any) *sql.Row
}

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

func (r *RecipeRepository) GetRecipeByID(ctx context.Context, id int64) (domain.Recipe, error) {
	const query = `
		SELECT
			r.id,
			r.name,
			i.id,
			i.name,
			i.unit,
			ri.quantity
		FROM recipes r
		LEFT JOIN recipe_ingredients ri ON ri.recipe_id = r.id
		LEFT JOIN ingredients i ON i.id = ri.ingredient_id
		WHERE r.id = $1
		ORDER BY i.id
	`

	rows, err := r.db.QueryContext(ctx, query, id)
	if err != nil {
		return domain.Recipe{}, fmt.Errorf("get recipe by id: %w", err)
	}
	defer rows.Close()

	var (
		recipe domain.Recipe
		found  bool
	)

	for rows.Next() {
		var (
			ingredientID sql.NullInt64
			ingredientName,
			unit sql.NullString
			quantity sql.NullInt64
		)

		if err := rows.Scan(
			&recipe.ID,
			&recipe.Name,
			&ingredientID,
			&ingredientName,
			&unit,
			&quantity,
		); err != nil {
			return domain.Recipe{}, fmt.Errorf("scan recipe by id: %w", err)
		}

		found = true

		if ingredientID.Valid {
			recipe.Ingredients = append(recipe.Ingredients, domain.RecipeIngredientLine{
				IngredientID: ingredientID.Int64,
				Name:         ingredientName.String,
				Unit:         unit.String,
				Quantity:     int(quantity.Int64),
			})
		}
	}

	if err := rows.Err(); err != nil {
		return domain.Recipe{}, fmt.Errorf("iterate recipe rows: %w", err)
	}

	if !found {
		return domain.Recipe{}, fmt.Errorf("recipe %d not found: %w", id, sql.ErrNoRows)
	}

	return recipe, nil
}

func (r *RecipeRepository) GetOrCreateIngredient(ctx context.Context, name, unit string) (int64, error) {
	id, err := getOrCreateIngredient(ctx, r.db, name, unit)
	if err != nil {
		return 0, fmt.Errorf("get or create ingredient: %w", err)
	}
	return id, nil
}

func getOrCreateIngredient(ctx context.Context, q queryRower, name, unit string) (int64, error) {
	const query = `
		INSERT INTO ingredients (name, unit)
		VALUES ($1, $2)
		ON CONFLICT (name) DO UPDATE SET unit = EXCLUDED.unit
		RETURNING id
	`

	var id int64
	if err := q.QueryRowContext(ctx, query, name, unit).Scan(&id); err != nil {
		return 0, err
	}

	return id, nil
}

func (r *RecipeRepository) SaveRecipe(ctx context.Context, recipe domain.Recipe) error {
	tx, err := r.db.BeginTx(ctx, nil)
	if err != nil {
		return fmt.Errorf("begin transaction: %w", err)
	}
	defer tx.Rollback()

	var recipeID int64
	const insertRecipe = `
		INSERT INTO recipes (name)
		VALUES ($1)
		RETURNING id
	`
	if err := tx.QueryRowContext(ctx, insertRecipe, recipe.Name).Scan(&recipeID); err != nil {
		return fmt.Errorf("insert recipe: %w", err)
	}

	const linkIngredient = `
		INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity)
		VALUES ($1, $2, $3)
	`

	for _, line := range recipe.Ingredients {
		ingredientID, err := getOrCreateIngredient(ctx, tx, line.Name, line.Unit)
		if err != nil {
			return fmt.Errorf("ingredient %q: %w", line.Name, err)
		}

		if _, err := tx.ExecContext(ctx, linkIngredient, recipeID, ingredientID, line.Quantity); err != nil {
			return fmt.Errorf("link ingredient %q: %w", line.Name, err)
		}
	}

	if err := tx.Commit(); err != nil {
		return fmt.Errorf("commit transaction: %w", err)
	}

	return nil
}