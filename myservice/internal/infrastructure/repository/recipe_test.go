package repository_test

import (
	"context"
	"os"
	"strings"
	"testing"

	"myservice/internal/config"
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