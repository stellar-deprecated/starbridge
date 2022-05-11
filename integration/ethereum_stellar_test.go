package integration

import (
	"io/ioutil"
	"net/http"
	"testing"
	"time"

	"github.com/stellar/go/keypair"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/require"
)

func TestEthereumStellarDeposit(t *testing.T) {
	itest := NewTest(t)

	keys, accounts := itest.CreateAccounts(1)
	key, account := keys[0], accounts[0]

	clientKey := keypair.MustParseFull("SBEICGMVMPF2WWIYV34IP7ON2Q6BUOT7F7IGHOTUMYUIG5K4IWIOUQC3")

	// Configure Starbridge account
	ops := []txnbuild.Operation{
		&txnbuild.CreateAccount{
			SourceAccount: account.GetAccountID(),
			Destination:   "GATBFH6GV7GMWNI5RXH546BB2MDSNO3DPLGPT4EAFS5ICLRZT3D7F4YS",
			Amount:        "100",
		},
		&txnbuild.ChangeTrust{
			SourceAccount: "GATBFH6GV7GMWNI5RXH546BB2MDSNO3DPLGPT4EAFS5ICLRZT3D7F4YS",
			Line: txnbuild.ChangeTrustAssetWrapper{
				Asset: txnbuild.CreditAsset{
					Code:   "ETH",
					Issuer: account.GetAccountID(),
				},
			},
		},
	}

	itest.MustSubmitMultiSigOperations(account, []*keypair.Full{key, clientKey}, ops...)

	var txEnvelope string

loop:
	for {
		time.Sleep(time.Second)
		resp, err := itest.Client().Get("http://localhost:8001/stellar/get_inverse_transaction/ethereum")
		require.NoError(t, err)
		switch resp.StatusCode {
		case http.StatusAccepted:
			t.Log("Signing request accepted, waiting...")
			continue loop
		case http.StatusOK:
			t.Log("Signing request processed")
		default:
			require.Failf(t, "Unknown status code", "Response code: %s", resp.Status)
		}
		body, err := ioutil.ReadAll(resp.Body)
		require.NoError(t, err)
		txEnvelope = string(body)
		break
	}

	t.Log(txEnvelope)
}
