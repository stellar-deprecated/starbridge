package main

import (
	"github.com/stellar/go/support/log"
	"github.com/stellar/starbridge/cmd"
)

func main() {
	err := cmd.RootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
