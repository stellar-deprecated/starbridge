package integrations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMyContractAddress(t *testing.T) {
	assert.True(t, IsMyContractAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48"), "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB48")
	assert.False(t, IsMyContractAddress("0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB00"), "0xA0b86991c6218b36c1d19D4a2e9Eb0cE3606eB00")
}
