package integration

import (
	"context"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"strconv"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/stellar/go/strkey"
	"github.com/stellar/starbridge/solidity-go"

	"github.com/stellar/starbridge/store"

	"github.com/stellar/go/support/errors"
	"github.com/stellar/go/txnbuild"
	"github.com/stretchr/testify/require"
	"golang.org/x/sync/errgroup"
)

const (
	ethereumSenderPrivateKey = "c1a4af60400ffd1473ada8425cff9f91b533194d6dd30424a17f356e418ac35b"
)

func depositETHToBridge(t *testing.T, client *ethclient.Client, amount *big.Int, stellarRecipient string) *types.Receipt {
	parsedPrivateKey, err := crypto.HexToECDSA(ethereumSenderPrivateKey)
	require.NoError(t, err)

	opts, err := bind.NewKeyedTransactorWithChainID(parsedPrivateKey, big.NewInt(31337))
	require.NoError(t, err)
	opts.Value = amount
	opts.GasPrice = new(big.Int).Mul(big.NewInt(1), big.NewInt(params.GWei))

	rawRecipient := strkey.MustDecode(strkey.VersionByteAccountID, stellarRecipient)
	recipient := &big.Int{}
	recipient.SetBytes(rawRecipient)

	bridge, err := solidity.NewBridgeTransactor(common.HexToAddress(EthereumBridgeAddress), client)
	require.NoError(t, err)
	tx, err := bridge.DepositETH(opts, recipient)
	require.NoError(t, err)

	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	require.Len(t, receipt.Logs, 1)
	return receipt
}

func TestEthereumToStellarWithdrawal(t *testing.T) {
	const servers int = 3

	itest := NewIntegrationTest(t, Config{
		Servers: servers,
	})

	ethRPCClient, err := ethclient.Dial(EthereumRPCURL)
	require.NoError(t, err)

	txs := make([]string, servers)
	stores := make([]*store.DB, servers)
	g := new(errgroup.Group)

	receipt := depositETHToBridge(
		t,
		ethRPCClient,
		new(big.Int).Mul(big.NewInt(3), big.NewInt(params.Ether)),
		itest.clientKey.Address(),
	)

	postData := url.Values{
		"transaction_hash": {receipt.TxHash.String()},
		"log_index":        {strconv.FormatUint(uint64(receipt.Logs[0].Index), 10)},
	}

	for i := 0; i < servers; i++ {
		i := i
		stores[i] = itest.app[i].NewStore()
		g.Go(func() error {
			port := 9000 + i
		loop:
			for {
				time.Sleep(time.Second)
				url := fmt.Sprintf("http://localhost:%d/stellar/withdraw/ethereum", port)
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
		require.Equal(t, expectedSeqNum, tx.SequenceNumber())

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

	memoHex := mainTx.ToXDR().Memo().MustHash().HexString()
	numFound := 0
	for attempts := 0; attempts < 10; attempts++ {
		for i := 0; i < servers; i++ {
			found, err := stores[i].HistoryStellarTransactionExists(context.Background(), memoHex)
			require.NoError(t, err)
			if found {
				numFound++
			}
		}
		if numFound == servers {
			break
		} else {
			numFound = 0
			time.Sleep(time.Second)
		}
	}
	require.Equal(t, servers, numFound)

	for i := 0; i < servers; i++ {
		port := 9000 + i
		time.Sleep(time.Second)
		url := fmt.Sprintf("http://localhost:%d/stellar/withdraw/ethereum", port)
		resp, err := itest.Client().PostForm(url, postData)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}
}
