package repository

import (
	"fmt"
	"os"
	"testing"

	"github.com/caarlos0/env"
	"github.com/distuurbia/PriceGenerator/internal/config"
	"github.com/distuurbia/PriceGenerator/internal/model"
	"github.com/go-redis/redis/v8"
	"github.com/ory/dockertest"
	"github.com/sirupsen/logrus"
)

var (
	priceGeneratorRepo *PriceGeneratorRepository
	shares             []*model.Share
)

func SetupRedis() (*redis.Client, func(), error) {
	pool, err := dockertest.NewPool("")
	if err != nil {
		return nil, nil, fmt.Errorf("could not construct pool: %w", err)
	}
	resource, err := pool.Run("redis", "latest", []string{})
	if err != nil {
		return nil, nil, fmt.Errorf("could not run the pool: %w", err)
	}

	client := redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("localhost:%s", resource.GetPort("6379/tcp")),
		DB:   0,
	})
	cleanup := func() {
		client.Close()
		pool.Purge(resource)
	}
	return client, cleanup, nil
}

func TestMain(m *testing.M) {
	const (
		berkshirePrice = 453066
		twitterPrice   = 54
		teslaPrice     = 290
		applePrice     = 193
		cocaColaPrice  = 62
	)
	shares = append(shares,
		&model.Share{Name: "Berkshire Hathaway Inc.", Price: berkshirePrice},
		&model.Share{Name: "Twitter", Price: twitterPrice},
		&model.Share{Name: "Tesla", Price: teslaPrice},
		&model.Share{Name: "Apple", Price: applePrice},
		&model.Share{Name: "Coca-Cola", Price: cocaColaPrice})

	rdsClient, cleanupRds, err := SetupRedis()
	if err != nil {
		fmt.Println(err)
		cleanupRds()
		os.Exit(1)
	}
	var cfg config.Config
	if err := env.Parse(&cfg); err != nil {
		logrus.Fatalf("failed to parse config: %v", err)
	}
	priceGeneratorRepo = NewPriceGeneratorRepository(rdsClient, &cfg)

	exitCode := m.Run()

	cleanupRds()
	os.Exit(exitCode)
}
