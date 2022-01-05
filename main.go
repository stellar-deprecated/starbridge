package main

import (
	"fmt"
	"log"

	"github.com/stellar/starbridge/transform"
)

func main() {
	txHash := "0x13070f64d40f22cd10c5bf9972767b67406ed3d818a50f82b1409289dcaa1aec"

	tx, e := transform.FetchEthTxByHash(txHash)
	if e != nil {
		log.Fatal(fmt.Errorf("error doing FetchEthTxByHash: %s", e))
	}
	fmt.Println(tx.String())

	stellarTx, e := transform.Transaction2Stellar(tx)
	if e != nil {
		log.Fatal(fmt.Errorf("error doing Transaction2Stellar: %s", e))
	}
	fmt.Println(transform.Stellar2String(stellarTx))
}
