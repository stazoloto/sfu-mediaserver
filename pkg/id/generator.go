package id

import (
	"crypto/rand"
	"fmt"
	"math/big"
)

// Generate генерирует новые ID
func Generate(minID, maxID int64) (int64, error) {

	if minID >= maxID {
		return 0, fmt.Errorf("invalid range: min >= max")
	}
	rangeVal := big.NewInt(maxID - minID + 1)

	randomNum, err := rand.Int(rand.Reader, rangeVal)
	if err != nil {
		return 0, fmt.Errorf("failed to generate random num: %w", err)
	}

	id := minID + randomNum.Int64()

	return id, nil
}
