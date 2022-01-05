package main

import (
	"fmt"
	"log"

	"github.com/stellar/starbridge/model"
	"github.com/stellar/starbridge/transform"
)

func main() {
	txHash := "0x13070f64d40f22cd10c5bf9972767b67406ed3d818a50f82b1409289dcaa1aec"

	modelTxEth, e := transform.FetchEthTxByHash(txHash)
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

	stellarTx, e := transform.Transaction2Stellar(modelTxStellar)
	if e != nil {
		log.Fatal(fmt.Errorf("error doing Transaction2Stellar: %s", e))
	}
	fmt.Println(transform.Stellar2String(stellarTx))
}
