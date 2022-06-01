package mjml

import (
	"crypto/rand"
	"fmt"
	"math"
	"math/big"
)

func randomIdentifier() (int32, error) {
	// generate a random number between 0 and the largest possible int32
	num, err := rand.Int(rand.Reader, big.NewInt(math.MaxInt32))
	if err != nil {
		return -1, fmt.Errorf("failed to generate random int: %w", err)
	}

	return int32(num.Int64()), nil
}
