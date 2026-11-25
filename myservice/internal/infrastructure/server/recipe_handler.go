package server

import (
	"context"
	"database/sql"
	"errors"

	"go.uber.org/zap"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"myservice/api"
	"myservice/internal/domain"
	"myservice/internal/infrastructure/repository"
)

type RecipeGRPCHandler struct {
	api.UnimplementedRecipeServiceServer

	repo *repository.RecipeRepository
	log  *zap.Logger
}

func NewRecipeGRPCHandler(repo *repository.RecipeRepository, log *zap.Logger) *RecipeGRPCHandler {
	return &RecipeGRPCHandler{repo: repo, log: log}
}

func (h *RecipeGRPCHandler) GetRecipe(ctx context.Context, req *api.GetRecipeRequest) (*api.GetRecipeResponse, error) {
	recipe, err := h.repo.GetRecipeByID(ctx, req.GetId())
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, status.Errorf(codes.NotFound, "recipe %d not found", req.GetId())
		}
		if h.log != nil {
			h.log.Error("get recipe by id failed",
				zap.Int64("id", req.GetId()),
				zap.Error(err),
			)
		}
		return nil, status.Errorf(codes.Internal, "get recipe: %v", err)
	}

	return &api.GetRecipeResponse{Recipe: toProtoRecipe(recipe)}, nil
}

func toProtoRecipe(r domain.Recipe) *api.Recipe {
	out := &api.Recipe{
		Id:   r.ID,
		Name: r.Name,
	}
	for _, line := range r.Ingredients {
		out.Ingredients = append(out.Ingredients, &api.Ingredient{
			Id:       line.IngredientID,
			Name:     line.Name,
			Unit:     line.Unit,
			Quantity: int32(line.Quantity),
		})
	}
	return out
}