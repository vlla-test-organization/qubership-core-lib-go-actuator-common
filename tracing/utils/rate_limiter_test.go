package utils

import (
	"github.com/stretchr/testify/assert"
	"math"
	"testing"
)

func TestNewRateLimiter(t *testing.T) {
	var rateLimiter *ReconfigurableRateLimiter
	var maxTracesPerSecond = float64(10)
	rateLimiter = NewRateLimiter(maxTracesPerSecond, math.Max(maxTracesPerSecond, 1.0))
	assert.Equal(t, float64(10), rateLimiter.creditsPerSecond)
	assert.Equal(t, float64(10), rateLimiter.maxBalance)
	assert.Equal(t, true, rateLimiter.CheckCredit(float64(10)))
	assert.Equal(t, false, rateLimiter.CheckCredit(float64(11)))
}

func TestCheckCredit(t *testing.T) {
	var rateLimiter *ReconfigurableRateLimiter
	var maxTracesPerSecond = float64(10)
	rateLimiter = NewRateLimiter(maxTracesPerSecond, math.Max(maxTracesPerSecond, 1.0))
	assert.Equal(t, true, rateLimiter.CheckCredit(float64(10)))
	assert.Equal(t, false, rateLimiter.CheckCredit(float64(11)))
}