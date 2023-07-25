package repository

import (
	"context"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestAddToStream(t *testing.T) {
	err := priceGeneratorRepo.AddToStream(context.Background(), shares)
	require.NoError(t, err)
}
