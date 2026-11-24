package kafka

import (
	"context"

	"go.uber.org/zap"

	"myservice/internal/application/processor"
)

type Consumer struct {
	log *zap.Logger
}

func NewConsumer(
	_ []string,
	_, _ string,
	_ *processor.RecipeProcessor,
	log *zap.Logger,
) (*Consumer, error) {
	return &Consumer{log: log}, nil
}

func (c *Consumer) Run(ctx context.Context) error {
	<-ctx.Done()
	return ctx.Err()
}

func (c *Consumer) Shutdown() {}