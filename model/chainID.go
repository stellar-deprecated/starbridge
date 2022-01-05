package model

import "fmt"

// Enum for ChainID
type ChainID uint8

const (
	ChainStellar ChainID = iota
	ChainEthereum
)

// String is the Stringer method
func (c ChainID) String() string {
	if c == ChainStellar {
		return "Stellar"
	} else if c == ChainEthereum {
		return "Ethereum"
	}

	// panic and display with %d because ChainID is a wrapper around a uint8; displaying as %v will be unintentionally recursive
	panic(fmt.Errorf("no such chain: %d", c))
}
