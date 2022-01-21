package main

import (
	"fmt"
	"log"

	"github.com/stellar/starbridge/integrations"
	"github.com/stellar/starbridge/model"
	"github.com/stellar/starbridge/transform"
)

func main() {
	txHash := "0x9a5ed1a2f961cbe3ddbf9ec083f662f0948924368bb8ea232b8abc5e1bfa70da"

	modelTxEth, e := integrations.FetchEthTxByHash(txHash)
	if e != nil {
		log.Fatal(fmt.Errorf("error doing FetchEthTxByHash: %s", e))
	}
	fmt.Println("transaction fetched as modelTxEth:")
	fmt.Println(modelTxEth.String())
	fmt.Printf("\n\n")

	modelTxStellar, e := transform.MapTxToChain(modelTxEth, model.ChainStellar)
	if e != nil {
		log.Fatal(fmt.Errorf("error doing MapTxToChain: %s", e))
	}
	fmt.Println("transaction converted to modelTxStellar:")
	fmt.Println(modelTxStellar.String())
	fmt.Printf("\n\n")

	stellarTx, e := integrations.Transaction2Stellar(modelTxStellar)
	if e != nil {
		log.Fatal(fmt.Errorf("error doing Transaction2Stellar: %s", e))
	}
	fmt.Println(integrations.Stellar2String(stellarTx))
}
