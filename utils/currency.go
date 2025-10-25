package utils

import (
	"math/rand"
	"time"
)


// ComputeEstimatedGDP estimates GDP given population and exchange rate.
func ComputeEstimatedGDP(population int64, exchangeRate float64) float64 {
	var rng = rand.New(rand.NewSource(time.Now().UnixNano()))

	randomNumber := rng.Intn(1001) + 1000

	estimatedGDP := float64(population) * float64(randomNumber) / exchangeRate

	return estimatedGDP
}
