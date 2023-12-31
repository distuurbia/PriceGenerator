// Package repository contains methods that work with Redis Stream
package repository

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/distuurbia/PriceGenerator/internal/config"
	"github.com/distuurbia/PriceGenerator/internal/model"
	"github.com/go-redis/redis/v8"
)

// PriceGeneratorRepository contains redis client
type PriceGeneratorRepository struct {
	client *redis.Client
	cfg    *config.Config
}

// NewPriceGeneratorRepository creates and returns a new instance of PriceGeneratorRepository, using the provided redis.Client
func NewPriceGeneratorRepository(client *redis.Client, cfg *config.Config) *PriceGeneratorRepository {
	return &PriceGeneratorRepository{
		client: client,
		cfg:    cfg,
	}
}

// AddToStream adds a message to the redis stream
func (r *PriceGeneratorRepository) AddToStream(ctx context.Context, shares []*model.Share) error {
	sharesJSON, err := json.Marshal(shares)
	if err != nil {
		return fmt.Errorf("RedisStreamRepository -> AddToStream -> json.Marshal -> %w", err)
	}
	streamData := redis.XAddArgs{
		Stream: r.cfg.RedisStreamName,
		Values: map[string]interface{}{
			r.cfg.RedisStreamField: string(sharesJSON),
		},
	}
	_, err = r.client.XAdd(ctx, &streamData).Result()
	if err != nil {
		return fmt.Errorf("RedisStreamRepository -> AddToStream -> XAdd -> %w", err)
	}
	return nil
}
