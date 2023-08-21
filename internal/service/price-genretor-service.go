// Package service contains bisnes logic of Price Generator
package service

import (
	"context"
	"crypto/rand"
	"fmt"
	"math/big"

	"github.com/distuurbia/PriceGenerator/internal/model"
)

const (
	minPriceChange = -1
	maxPriceChange = 1
	intSize        = 53
)

// PriceGeneratorRepository is an interface of PriceGeneratorRepository structure of repository
type PriceGeneratorRepository interface {
	AddToStream(ctx context.Context, shares []*model.Share) error
}

// PriceGeneratorService contains an inerface of PriceGeneratorRepository
type PriceGeneratorService struct {
	r PriceGeneratorRepository
}

// NewPriceGeneratorService creates an object of PriceGeneratorService by using PriceGeneratorRepository interface
func NewPriceGeneratorService(r PriceGeneratorRepository) *PriceGeneratorService {
	return &PriceGeneratorService{r: r}
}

// GenerateRandomDeltaPrice returns random delta for changing price
func (s *PriceGeneratorService) GenerateRandomDeltaPrice() (delta float64, err error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(1<<intSize))
	if err != nil {
		return delta, fmt.Errorf("GenereteRandomDeltaPrice %w ", err)
	}
	delta = (float64(nBig.Int64())/(1<<intSize))*(maxPriceChange-minPriceChange) + minPriceChange
	return delta, nil
}

// SharesChangePrice return shares with changed price by some random delta
func (s *PriceGeneratorService) SharesChangePrice(shares []*model.Share) (sharesChangedPrice []*model.Share, err error) {
	for _, share := range shares {
		delta, err := s.GenerateRandomDeltaPrice()
		if err != nil {
			return nil, fmt.Errorf("SharesChangePrice -> %w ", err)
		}
		share.Price += delta
		if share.Price <= 0 {
			share.Price++
		}
	}
	return shares, nil
}

// AddToStream adds to stream array of shares trough repository AddToStream method
func (s *PriceGeneratorService) AddToStream(ctx context.Context, shares []*model.Share) (newShares []*model.Share, err error) {
	shares, err = s.SharesChangePrice(shares)
	if err != nil {
		return nil, fmt.Errorf("PriceGeneratorService -> AddToStream -> %w ", err)
	}
	err = s.r.AddToStream(ctx, shares)
	if err != nil {
		return nil, fmt.Errorf("PriceGeneratorService -> AddToStream -> %w", err)
	}
	return shares, nil
}
