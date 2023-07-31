// Package main contains main func and redis connection
package main

import (
	"context"
	"time"

	"github.com/caarlos0/env/v8"
	"github.com/distuurbia/PriceGenerator/internal/config"
	"github.com/distuurbia/PriceGenerator/internal/model"
	"github.com/distuurbia/PriceGenerator/internal/repository"
	"github.com/distuurbia/PriceGenerator/internal/service"
	"github.com/go-redis/redis/v8"
	"github.com/sirupsen/logrus"
)

const (
	berkshirePrice = 453066
	twitterPrice   = 54
	teslaPrice     = 290
	applePrice     = 193
	cocaColaPrice  = 62
)

// connectRedis connects to the redis db
func connectRedis(cfg *config.Config) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: cfg.RedisAddress,
		DB:   0,
	})
	return client
}

// createShares creates some hardcoded shares slice
func createShares() []*model.Share {
	var shares []*model.Share
	shares = append(shares,
		&model.Share{Name: "Berkshire Hathaway Inc.", Price: berkshirePrice},
		&model.Share{Name: "Twitter", Price: twitterPrice},
		&model.Share{Name: "Tesla", Price: teslaPrice},
		&model.Share{Name: "Apple", Price: applePrice},
		&model.Share{Name: "Coca-Cola", Price: cocaColaPrice})
	return shares
}

func main() {
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		logrus.Fatalf("failed to parse config: %v", err)
	}

	client := connectRedis(&cfg)
	r := repository.NewPriceGeneratorRepository(client, &cfg)
	s := service.NewPriceGeneratorService(r)
	shares := createShares()
	var err error
	for {
		shares, err = s.AddToStream(context.Background(), shares)
		if err != nil {
			logrus.Fatalf("main -> %v", err)
		}
		time.Sleep(1 * time.Second)
	}
}
