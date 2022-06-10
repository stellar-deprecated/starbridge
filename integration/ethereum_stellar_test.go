package integration

import (
	"context"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"testing"
	"time"

	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/starbridge/store"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestEthereumStellarDeposit(t *testing.T) {
	const servers int = 3

	itest := NewIntegrationTest(t, Config{
		Servers: servers,
	})

	incomingTx := store.IncomingEthereumTransaction{
		Hash:           "bf308af417b896b78f1a6bc5b8bd53df1a6d0270ba17c64345dac01b21d9559f",
		ValueWei:       1000,
		StellarAddress: itest.clientKey.Address(),
	}

	for i := 0; i < servers; i++ {
		err := itest.app[i].GetStore().InsertIncomingEthereumTransaction(context.Background(), incomingTx)
		require.NoError(t, err)
		_, err = itest.app[i].GetStore().GetIncomingEthereumTransactionByHash(context.Background(), incomingTx.Hash)
		require.NoError(t, err)
	}

	txs := make([]string, servers)

	g := new(errgroup.Group)

	postData := url.Values{
		"transaction_hash": {incomingTx.Hash},
	}

	for i := 0; i < servers; i++ {
		i := i
		g.Go(func() error {
			port := 9000 + i
		loop:
			for {
				time.Sleep(time.Second)
				url := fmt.Sprintf("http://localhost:%d/stellar/get_inverse_transaction/ethereum", port)
				resp, err := itest.Client().PostForm(url, postData)
				require.NoError(t, err)
				switch resp.StatusCode {
				case http.StatusAccepted:
					t.Log("Signing request accepted, waiting...")
					continue loop
				case http.StatusOK:
					t.Log("Signing request processed")
				default:
					return errors.Errorf("Unknown status code: %s", resp.Status)
				}
				body, err := ioutil.ReadAll(resp.Body)
				require.NoError(t, err)
				txEnvelope := string(body)
				txs[i] = txEnvelope

				// Try to submit tx with just one signature
				_, err = itest.HorizonClient().SubmitTransactionXDR(txEnvelope)
				require.Error(t, err)

				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}

	// Concat signatures
	gtx, err := txnbuild.TransactionFromXDR(txs[0])
	require.NoError(t, err)
	mainTx, ok := gtx.Transaction()
	require.True(t, ok)

	expectedSeqNum := mainTx.SequenceNumber()

	// Add as many sigs as needed and not a single more
	for i := 1; i < servers/2+1; i++ {
		gtx, err := txnbuild.TransactionFromXDR(txs[i])
		require.NoError(t, err)
		tx, ok := gtx.Transaction()
		require.True(t, ok)

		// TODO timebound time should be provided in HTTP request and checked by Starbridge
		assert.Equal(t, expectedSeqNum, tx.SequenceNumber())

		sig := tx.Signatures()

		mainTx, err = mainTx.AddSignatureDecorated(sig...)
		require.NoError(t, err)
	}

	// ...and add client signature (because it's tx source)
	mainTx, err = mainTx.Sign(StandaloneNetworkPassphrase, itest.clientKey)
	require.NoError(t, err)

	b64, err := mainTx.Base64()
	require.NoError(t, err)

	_, err = itest.HorizonClient().SubmitTransaction(mainTx)
	require.NoErrorf(t, err, "error submitting: %s", b64)
}
