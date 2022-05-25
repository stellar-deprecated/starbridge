package integration

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

func TestEthereumStellarDeposit(t *testing.T) {
	const servers int = 3

	itest := NewIntegrationTest(t, Config{
		Servers: servers,
	})

	txs := make([]string, servers)

	g := new(errgroup.Group)
	for i := 0; i < servers; i++ {
		i := i
		g.Go(func() error {
			port := 9000 + i
		loop:
			for {
				time.Sleep(time.Second)
				postData := url.Values{
					"tx_expiration_timestamp": {strconv.FormatInt(time.Now().Add(time.Minute).Unix(), 10)},
				}
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

	_, err = itest.HorizonClient().SubmitTransaction(mainTx)
	require.NoError(t, err)
}
