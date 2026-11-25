package repository_test

import (
	"context"
	"database/sql"
	"errors"
	"os"
	"strings"
	"testing"

	"myservice/internal/config"
	"myservice/internal/domain"
	"myservice/internal/infrastructure/database"
	"myservice/internal/infrastructure/repository"
)

func loadDotEnv(paths ...string) {
	for _, path := range paths {
		data, err := os.ReadFile(path)
		if err != nil {
			continue
		}
		for line := range strings.SplitSeq(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			key, value, ok := strings.Cut(line, "=")
			if !ok {
				continue
			}
			key = strings.TrimSpace(key)
			value = strings.TrimSpace(value)
			if os.Getenv(key) == "" {
				os.Setenv(key, value)
			}
		}
		return
	}
}

func testDBConfig(t *testing.T) config.Config {
	t.Helper()

	loadDotEnv(".env", "../../.env", "../../../.env")

	cfg, err := config.LoadConfig()
	if err != nil {
		t.Fatalf("load config: %v", err)
	}

	cfg.DBName = os.Getenv("TEST_DB_NAME")
	if cfg.DBName == "" {
		cfg.DBName = "mixturka_test"
	}

	return cfg
}

func TestCreateRecipe(t *testing.T) {
	ctx := context.Background()
	cfg := testDBConfig(t)

	db, err := database.NewPostgres(cfg)
	if err != nil {
		t.Skipf("postgres unavailable (docker + test DB %q?): %v", cfg.DBName, err)
	}
	t.Cleanup(func() { _ = db.Close() })

	repo := repository.NewRecipeRepository(db)

	const wantName = "Тестовый борщ"
	id, err := repo.CreateRecipe(ctx, wantName)
	if err != nil {
		t.Fatalf("CreateRecipe: %v", err)
	}
	if id <= 0 {
		t.Fatalf("id = %d, want > 0", id)
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, "DELETE FROM recipes WHERE id = $1", id)
	})

	var gotName string
	err = db.QueryRowContext(ctx, `SELECT name FROM recipes WHERE id = $1`, id).Scan(&gotName)
	if err != nil {
		t.Fatalf("select recipe by id: %v", err)
	}
	if gotName != wantName {
		t.Fatalf("name = %q, want %q", gotName, wantName)
	}
}

func openTestDB(t *testing.T) (*sql.DB, *repository.RecipeRepository) {
	t.Helper()

	cfg := testDBConfig(t)
	db, err := database.NewPostgres(cfg)
	if err != nil {
		t.Skipf("postgres unavailable (docker + test DB %q?): %v", cfg.DBName, err)
	}
	t.Cleanup(func() { _ = db.Close() })

	return db, repository.NewRecipeRepository(db)
}

func seedRecipeWithIngredients(t *testing.T, db *sql.DB, ctx context.Context) int64 {
	t.Helper()

	var recipeID int64
	err := db.QueryRowContext(ctx,
		`INSERT INTO recipes (name) VALUES ($1) RETURNING id`, "Борщ",
	).Scan(&recipeID)
	if err != nil {
		t.Fatalf("insert recipe: %v", err)
	}

	type ing struct {
		name string
		unit string
		qty  int
	}
	items := []ing{
		{"Свёкла", "grams", 200},
		{"Картофель", "grams", 400},
		{"Морковь", "grams", 100},
	}

	for _, item := range items {
		var ingredientID int64
		err := db.QueryRowContext(ctx, `
			INSERT INTO ingredients (name, unit) VALUES ($1, $2)
			ON CONFLICT (name) DO UPDATE SET unit = EXCLUDED.unit
			RETURNING id`,
			item.name, item.unit,
		).Scan(&ingredientID)
		if err != nil {
			t.Fatalf("insert ingredient %q: %v", item.name, err)
		}

		_, err = db.ExecContext(ctx, `
			INSERT INTO recipe_ingredients (recipe_id, ingredient_id, quantity)
			VALUES ($1, $2, $3)`,
			recipeID, ingredientID, item.qty,
		)
		if err != nil {
			t.Fatalf("link ingredient %q: %v", item.name, err)
		}
	}

	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM recipes WHERE id = $1`, recipeID)
	})

	return recipeID
}

func TestGetRecipeByID(t *testing.T) {
	ctx := context.Background()
	db, repo := openTestDB(t)

	t.Run("three ingredients", func(t *testing.T) {
		id := seedRecipeWithIngredients(t, db, ctx)

		recipe, err := repo.GetRecipeByID(ctx, id)
		if err != nil {
			t.Fatalf("GetRecipeByID: %v", err)
		}
		if recipe.ID != id || recipe.Name != "Борщ" {
			t.Fatalf("recipe = %+v, want id=%d name=Борщ", recipe, id)
		}
		if len(recipe.Ingredients) != 3 {
			t.Fatalf("ingredients len = %d, want 3", len(recipe.Ingredients))
		}

		want := map[string]struct{ unit string; qty int }{
			"Свёкла":    {"grams", 200},
			"Картофель": {"grams", 400},
			"Морковь":   {"grams", 100},
		}
		for _, line := range recipe.Ingredients {
			spec, ok := want[line.Name]
			if !ok {
				t.Fatalf("unexpected ingredient %q", line.Name)
			}
			if line.Unit != spec.unit || line.Quantity != spec.qty {
				t.Fatalf("%q: unit=%q qty=%d, want %q %d", line.Name, line.Unit, line.Quantity, spec.unit, spec.qty)
			}
			if line.IngredientID <= 0 {
				t.Fatalf("%q: ingredient id = %d", line.Name, line.IngredientID)
			}
		}
	})

	t.Run("no ingredients", func(t *testing.T) {
		id, err := repo.CreateRecipe(ctx, "Пустой салат")
		if err != nil {
			t.Fatalf("CreateRecipe: %v", err)
		}
		t.Cleanup(func() {
			_, _ = db.ExecContext(ctx, `DELETE FROM recipes WHERE id = $1`, id)
		})

		recipe, err := repo.GetRecipeByID(ctx, id)
		if err != nil {
			t.Fatalf("GetRecipeByID: %v", err)
		}
		if len(recipe.Ingredients) != 0 {
			t.Fatalf("ingredients = %+v, want empty", recipe.Ingredients)
		}
	})

	t.Run("not found", func(t *testing.T) {
		_, err := repo.GetRecipeByID(ctx, 999999999)
		if err == nil {
			t.Fatal("expected error for missing recipe")
		}
		if !errors.Is(err, sql.ErrNoRows) {
			t.Fatalf("err = %v, want sql.ErrNoRows in chain", err)
		}
	})
}

func TestGetOrCreateIngredient(t *testing.T) {
	ctx := context.Background()
	db, repo := openTestDB(t)

	const seedName = "Картофель-TestGetOrCreateIngredient"
	const seedUnit = "grams"

	seededID, err := repo.GetOrCreateIngredient(ctx, seedName, seedUnit)
	if err != nil {
		t.Fatalf("seed ingredient: %v", err)
	}
	t.Cleanup(func() {
		_, _ = db.ExecContext(ctx, `DELETE FROM ingredients WHERE id = $1`, seededID)
	})

	t.Run("returns existing id", func(t *testing.T) {
		id, err := repo.GetOrCreateIngredient(ctx, seedName, seedUnit)
		if err != nil {
			t.Fatalf("GetOrCreateIngredient: %v", err)
		}
		if id != seededID {
			t.Fatalf("id = %d, want %d", id, seededID)
		}
	})

	t.Run("no duplicate on repeat", func(t *testing.T) {
		id1, err := repo.GetOrCreateIngredient(ctx, seedName, seedUnit)
		if err != nil {
			t.Fatalf("first call: %v", err)
		}
		id2, err := repo.GetOrCreateIngredient(ctx, seedName, seedUnit)
		if err != nil {
			t.Fatalf("second call: %v", err)
		}
		if id1 != id2 {
			t.Fatalf("ids differ: %d vs %d", id1, id2)
		}

		var count int
		err = db.QueryRowContext(ctx,
			`SELECT COUNT(*) FROM ingredients WHERE name = $1`, seedName,
		).Scan(&count)
		if err != nil {
			t.Fatalf("count: %v", err)
		}
		if count != 1 {
			t.Fatalf("count = %d, want 1", count)
		}
	})

	t.Run("creates new ingredient", func(t *testing.T) {
		const newName = "Морковь-get-or-create-test"
		const newUnit = "grams"

		id, err := repo.GetOrCreateIngredient(ctx, newName, newUnit)
		if err != nil {
			t.Fatalf("GetOrCreateIngredient: %v", err)
		}
		if id <= 0 {
			t.Fatalf("id = %d, want > 0", id)
		}
		if id == seededID {
			t.Fatalf("new ingredient got seeded id %d", seededID)
		}

		t.Cleanup(func() {
			_, _ = db.ExecContext(ctx, `DELETE FROM ingredients WHERE id = $1`, id)
		})

		var gotName, gotUnit string
		err = db.QueryRowContext(ctx,
			`SELECT name, unit FROM ingredients WHERE id = $1`, id,
		).Scan(&gotName, &gotUnit)
		if err != nil {
			t.Fatalf("select: %v", err)
		}
		if gotName != newName || gotUnit != newUnit {
			t.Fatalf("stored = %q/%q, want %q/%q", gotName, gotUnit, newName, newUnit)
		}
	})

	t.Run("postgres error", func(t *testing.T) {
		_, err := repo.GetOrCreateIngredient(ctx, "   ", "grams")
		if err == nil {
			t.Fatal("expected error for invalid name")
		}
	})
}

func TestSaveRecipe(t *testing.T) {
	ctx := context.Background()
	db, repo := openTestDB(t)

	t.Run("three ingredients", func(t *testing.T) {
		const recipeName = "Борщ-TestSaveRecipe"
		in := domain.Recipe{
			Name: recipeName,
			Ingredients: []domain.RecipeIngredientLine{
				{Name: "Свёкла-save", Unit: "grams", Quantity: 200},
				{Name: "Картофель-save", Unit: "grams", Quantity: 400},
				{Name: "Морковь-save", Unit: "grams", Quantity: 100},
			},
		}

		if err := repo.SaveRecipe(ctx, in); err != nil {
			t.Fatalf("SaveRecipe: %v", err)
		}

		var recipeID int64
		err := db.QueryRowContext(ctx, `SELECT id FROM recipes WHERE name = $1`, recipeName).Scan(&recipeID)
		if err != nil {
			t.Fatalf("select recipe id: %v", err)
		}
		t.Cleanup(func() {
			_, _ = db.ExecContext(ctx, `DELETE FROM recipes WHERE id = $1`, recipeID)
		})

		got, err := repo.GetRecipeByID(ctx, recipeID)
		if err != nil {
			t.Fatalf("GetRecipeByID: %v", err)
		}
		if got.Name != recipeName || len(got.Ingredients) != 3 {
			t.Fatalf("got = %+v", got)
		}

		wantQty := map[string]int{
			"Свёкла-save":    200,
			"Картофель-save": 400,
			"Морковь-save":   100,
		}
		for _, line := range got.Ingredients {
			q, ok := wantQty[line.Name]
			if !ok {
				t.Fatalf("unexpected ingredient %q", line.Name)
			}
			if line.Quantity != q || line.Unit != "grams" {
				t.Fatalf("%q: qty=%d unit=%q", line.Name, line.Quantity, line.Unit)
			}
		}
	})

	t.Run("reuses existing ingredient", func(t *testing.T) {
		const ingName = "Картофель-SaveRecipe-reuse"
		existingID, err := repo.GetOrCreateIngredient(ctx, ingName, "grams")
		if err != nil {
			t.Fatalf("seed: %v", err)
		}

		recipeName := "Суп-SaveRecipe-reuse"
		err = repo.SaveRecipe(ctx, domain.Recipe{
			Name: recipeName,
			Ingredients: []domain.RecipeIngredientLine{
				{Name: ingName, Unit: "grams", Quantity: 250},
			},
		})
		if err != nil {
			t.Fatalf("SaveRecipe: %v", err)
		}

		t.Cleanup(func() {
			_, _ = db.ExecContext(ctx, `DELETE FROM recipes WHERE name = $1`, recipeName)
		})

		var count int
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM ingredients WHERE name = $1`, ingName).Scan(&count)
		if err != nil {
			t.Fatalf("count: %v", err)
		}
		if count != 1 {
			t.Fatalf("ingredient count = %d, want 1", count)
		}

		var usedID int64
		err = db.QueryRowContext(ctx, `
			SELECT ingredient_id FROM recipe_ingredients ri
			JOIN recipes r ON r.id = ri.recipe_id
			WHERE r.name = $1`, recipeName,
		).Scan(&usedID)
		if err != nil {
			t.Fatalf("link: %v", err)
		}
		if usedID != existingID {
			t.Fatalf("used id %d, want %d", usedID, existingID)
		}
	})

	t.Run("rollback on error", func(t *testing.T) {
		const recipeName = "Bad-SaveRecipe-rollback"
		err := repo.SaveRecipe(ctx, domain.Recipe{
			Name: recipeName,
			Ingredients: []domain.RecipeIngredientLine{
				{Name: "Valid-ing-save", Unit: "grams", Quantity: 100},
				{Name: "Bad-qty-ing-save", Unit: "grams", Quantity: 0},
			},
		})
		if err == nil {
			t.Fatal("expected SaveRecipe error")
		}

		var recipeCount int
		err = db.QueryRowContext(ctx, `SELECT COUNT(*) FROM recipes WHERE name = $1`, recipeName).Scan(&recipeCount)
		if err != nil {
			t.Fatalf("count recipes: %v", err)
		}
		if recipeCount != 0 {
			t.Fatalf("recipes count = %d, want 0 after rollback", recipeCount)
		}

		var linkCount int
		err = db.QueryRowContext(ctx, `
			SELECT COUNT(*) FROM recipe_ingredients ri
			JOIN ingredients i ON i.id = ri.ingredient_id
			WHERE i.name IN ($1, $2)`,
			"Valid-ing-save", "Bad-qty-ing-save",
		).Scan(&linkCount)
		if err != nil {
			t.Fatalf("count links: %v", err)
		}
		if linkCount != 0 {
			t.Fatalf("link count = %d, want 0 after rollback", linkCount)
		}
	})
}