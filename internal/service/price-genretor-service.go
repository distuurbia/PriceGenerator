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
	minPriceChange = -0.1
	maxPriceChange = 0.1
	intSize        = 53
)

// PriceGeneratorRepository is an interface of PriceGeneratorRepository structure of repository
type PriceGeneratorRepository interface {
	AddToStream(ctx context.Context, shares []*model.Share) error
}

// PriceGeneratorService contains an inerface of PriceGeneratorRepository
type PriceGeneratorService struct {
	priceGeneratorRepo PriceGeneratorRepository
}

// NewPriceGeneratorService creates an object of PriceGeneratorService by using PriceGeneratorRepository interface
func NewPriceGeneratorService(priceGeneratorRepo PriceGeneratorRepository) *PriceGeneratorService {
	return &PriceGeneratorService{priceGeneratorRepo: priceGeneratorRepo}
}

// GenerateRandomDeltaPrice returns random delta for changing price
func (priceGeneratorSrv *PriceGeneratorService) GenerateRandomDeltaPrice() (delta float64, err error) {
	nBig, err := rand.Int(rand.Reader, big.NewInt(1<<intSize))
	if err != nil {
		return delta, fmt.Errorf("GenereteRandomDeltaPrice %w ", err)
	}
	delta = (float64(nBig.Int64())/(1<<intSize))*(maxPriceChange-minPriceChange) + minPriceChange
	return delta, nil
}

// SharesChangePrice return shares with changed price by some random delta
func (priceGeneratorSrv *PriceGeneratorService) SharesChangePrice(shares []*model.Share) (sharesChangedPrice []*model.Share, err error) {
	for _, share := range shares {
		delta, err := priceGeneratorSrv.GenerateRandomDeltaPrice()
		if err != nil {
			return nil, fmt.Errorf("SharesChangePrice -> %w ", err)
		}
		share.Price += delta
	}
	return shares, nil
}

// AddToStream adds to stream array of shares trough repository AddToStream method
func (priceGeneratorSrv *PriceGeneratorService) AddToStream(ctx context.Context, shares []*model.Share) (err error) {
	shares, err = priceGeneratorSrv.SharesChangePrice(shares)
	if err != nil {
		return fmt.Errorf("PriceGeneratorService -> AddToStream -> %w ", err)
	}
	err = priceGeneratorSrv.priceGeneratorRepo.AddToStream(ctx, shares)
	if err != nil {
		return fmt.Errorf("priceGeneratorSrv -> AddToStream -> %w", err)
	}
	return nil
}
