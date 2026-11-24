package processor

import (
	"go.uber.org/zap"

	"myservice/internal/infrastructure/repository"
)

type RecipeProcessor struct {
	repo *repository.RecipeRepository
	log  *zap.Logger
}

func NewRecipeProcessor(repo *repository.RecipeRepository, log *zap.Logger) *RecipeProcessor {
	return &RecipeProcessor{repo: repo, log: log}
}