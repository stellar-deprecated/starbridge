package integrations

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIsMyContractAddress(t *testing.T) {
	assert.True(t, IsMyContractAddress("0x35459293165463e68b17ce13d2ddd79654eae0d6"), "0x35459293165463e68b17ce13d2ddd79654eae0d6")
	assert.False(t, IsMyContractAddress("0x35459293165463e68b17ce13d2ddd79654eae0d0"), "0x35459293165463e68b17ce13d2ddd79654eae0d0")
}
