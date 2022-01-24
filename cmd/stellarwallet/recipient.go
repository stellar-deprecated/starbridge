package main

import (
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
)

func IsRecipient(tx *txnbuild.Transaction, account *keypair.FromAddress) bool {
	for _, op := range tx.Operations() {
		ccb, ok := op.(*txnbuild.CreateClaimableBalance)
		if !ok {
			continue
		}
		for _, c := range ccb.Destinations {
			// TODO: Assess the predicate portion of the claimant.
			if c.Destination == account.Address() {
				return true
			}
		}
	}
	return false
}
