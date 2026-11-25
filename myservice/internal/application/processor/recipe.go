package processor

import (
	"context"
	"encoding/json"
	"fmt"
	"strings"

	"go.uber.org/zap"

	"myservice/internal/application/processor/dto"
	"myservice/internal/domain"
	"myservice/internal/infrastructure/repository"
)

type RecipeProcessor struct {
	repo *repository.RecipeRepository
	log  *zap.Logger
}

func NewRecipeProcessor(repo *repository.RecipeRepository, log *zap.Logger) *RecipeProcessor {
	return &RecipeProcessor{repo: repo, log: log}
}

func (p *RecipeProcessor) ProcessMessage(ctx context.Context, data []byte) error {
	var msg dto.RecipeMessage
	if err := json.Unmarshal(data, &msg); err != nil {
		p.log.Error("invalid kafka message", zap.Error(err))
		return fmt.Errorf("unmarshal recipe message: %w", err)
	}

	return p.ProcessRecipe(ctx, msg)
}

func (p *RecipeProcessor) ProcessRecipe(ctx context.Context, msg dto.RecipeMessage) error {
	recipe, err := toDomainRecipe(msg)
	if err != nil {
		p.log.Warn("reject recipe message",
			zap.String("message_id", msg.ID),
			zap.String("recipe_name", msg.Name),
			zap.Error(err),
		)
		return err
	}

	if err := p.repo.SaveRecipe(ctx, recipe); err != nil {
		p.log.Error("save recipe failed",
			zap.String("message_id", msg.ID),
			zap.String("recipe_name", recipe.Name),
			zap.Error(err),
		)
		return fmt.Errorf("save recipe: %w", err)
	}

	p.log.Info("recipe processed successfully",
		zap.String("message_id", msg.ID),
		zap.String("recipe_name", recipe.Name),
		zap.Int("ingredients_count", len(recipe.Ingredients)),
	)

	return nil
}

func toDomainRecipe(msg dto.RecipeMessage) (domain.Recipe, error) {
	name := strings.TrimSpace(msg.Name)
	if name == "" {
		return domain.Recipe{}, fmt.Errorf("recipe name must not be empty")
	}

	recipe := domain.Recipe{Name: name}

	for i, ing := range msg.Ingredients {
		ingredientName := strings.TrimSpace(ing.Name)
		if ingredientName == "" {
			return domain.Recipe{}, fmt.Errorf("ingredient %d: name must not be empty", i)
		}
		if ing.Quantity <= 0 {
			return domain.Recipe{}, fmt.Errorf("ingredient %q: quantity must be greater than zero", ingredientName)
		}

		unit := strings.TrimSpace(ing.Unit)
		if unit == "" {
			return domain.Recipe{}, fmt.Errorf("ingredient %q: unit must not be empty", ingredientName)
		}

		recipe.Ingredients = append(recipe.Ingredients, domain.RecipeIngredientLine{
			Name:     ingredientName,
			Unit:     unit,
			Quantity: int(ing.Quantity),
		})
	}

	return recipe, nil
}
