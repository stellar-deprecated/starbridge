package integration

import (
	"context"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/big"
	"net/http"
	"net/url"
	"sort"
	"strconv"
	"testing"
	"time"

	"github.com/stellar/starbridge/stellar/controllers"

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
	require.Equal(t, uint64(1), receipt.Status)
	return receipt
}

func withdrawETHFromBridge(t *testing.T, client *ethclient.Client, numValidators int, amount *big.Int, responses []controllers.EthereumSignatureResponse) {
	parsedPrivateKey, err := crypto.HexToECDSA(ethereumSenderPrivateKey)
	require.NoError(t, err)

	opts, err := bind.NewKeyedTransactorWithChainID(parsedPrivateKey, big.NewInt(31337))
	require.NoError(t, err)
	opts.GasPrice = new(big.Int).Mul(big.NewInt(1), big.NewInt(params.GWei))

	bridge, err := solidity.NewBridgeTransactor(common.HexToAddress(EthereumBridgeAddress), client)
	require.NoError(t, err)

	caller, err := solidity.NewBridgeCaller(common.HexToAddress(EthereumBridgeAddress), client)
	require.NoError(t, err)
	validatorToIndex := map[common.Address]uint8{}
	for i := 0; i < numValidators; i++ {
		address, err := caller.Signers(nil, big.NewInt(int64(i)))
		require.NoError(t, err)
		validatorToIndex[address] = uint8(i)
	}

	sort.Slice(responses, func(i, j int) bool {
		address := common.HexToAddress(responses[i].Address)
		index, ok := validatorToIndex[address]
		require.True(t, ok)

		otherAddress := common.HexToAddress(responses[j].Address)
		otherIndex, ok := validatorToIndex[otherAddress]
		require.True(t, ok)

		return index < otherIndex
	})
	signatures := make([][]byte, len(responses))
	indexes := make([]uint8, len(responses))
	for i, response := range responses {
		signatures[i] = common.Hex2Bytes(response.Signature)
		indexes[i] = validatorToIndex[common.HexToAddress(response.Address)]
	}

	tx, err := bridge.WithdrawETH(
		opts,
		solidity.WithdrawETHRequest{
			Id:         common.HexToHash(responses[0].DepositID),
			Expiration: big.NewInt(responses[0].Expiration),
			Recipient:  opts.From,
			Amount:     amount,
		},
		signatures,
		indexes,
	)
	require.NoError(t, err)

	receipt, err := client.TransactionReceipt(context.Background(), tx.Hash())
	require.NoError(t, err)
	require.Equal(t, uint64(1), receipt.Status)
}

func TestEthereumToStellarWithdrawal(t *testing.T) {
	const servers int = 3

	itest := NewIntegrationTest(t, Config{
		Servers:                servers,
		EthereumFinalityBuffer: 0,
		WithdrawalWindow:       time.Hour,
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
				url := fmt.Sprintf("http://localhost:%d/ethereum/withdraw/stellar", port)
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

	// Too early for refunds
	for i := 0; i < servers; i++ {
		port := 9000 + i
		url := fmt.Sprintf("http://localhost:%d/ethereum/refund", port)
		resp, err := itest.Client().PostForm(url, postData)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
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

	// cannot withdraw more than once
	for i := 0; i < servers; i++ {
		port := 9000 + i
		time.Sleep(time.Second)
		url := fmt.Sprintf("http://localhost:%d/ethereum/withdraw/stellar", port)
		resp, err := itest.Client().PostForm(url, postData)
		require.NoError(t, err)
		require.Equal(t, http.StatusBadRequest, resp.StatusCode)
	}
}

func TestEthereumRefund(t *testing.T) {
	const servers int = 3

	itest := NewIntegrationTest(t, Config{
		Servers:                servers,
		EthereumFinalityBuffer: 0,
		WithdrawalWindow:       time.Second,
	})

	ethRPCClient, err := ethclient.Dial(EthereumRPCURL)
	require.NoError(t, err)

	responses := make([]controllers.EthereumSignatureResponse, servers)
	stores := make([]*store.DB, servers)
	for i := 0; i < servers; i++ {
		stores[i] = itest.app[i].NewStore()
	}

	depositAmount := new(big.Int).Mul(big.NewInt(3), big.NewInt(params.Ether))
	receipt := depositETHToBridge(
		t,
		ethRPCClient,
		depositAmount,
		itest.clientKey.Address(),
	)

	postData := url.Values{
		"transaction_hash": {receipt.TxHash.String()},
		"log_index":        {strconv.FormatUint(uint64(receipt.Logs[0].Index), 10)},
	}

	header, err := ethRPCClient.HeaderByHash(context.Background(), receipt.BlockHash)
	require.NoError(t, err)
	depositTime := time.Unix(int64(header.Time), 0)
	for {
		ready := 0
		for i := 0; i < servers; i++ {
			lastCloseTime, err := stores[i].GetLastLedgerCloseTime(context.Background())
			require.NoError(t, err)
			if lastCloseTime.After(depositTime.Add(time.Second)) {
				ready++
			}
		}
		if ready == servers {
			break
		}
		time.Sleep(time.Second)
	}

	g := new(errgroup.Group)

	for i := 0; i < servers; i++ {
		i := i
		stores[i] = itest.app[i].NewStore()
		g.Go(func() error {
			port := 9000 + i
		loop:
			for {
				time.Sleep(time.Second)
				url := fmt.Sprintf("http://localhost:%d/ethereum/refund", port)
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
				require.NoError(t, json.NewDecoder(resp.Body).Decode(&responses[i]))
				return nil
			}
		})
	}

	if err := g.Wait(); err != nil {
		t.Fatal(err)
	}

	withdrawETHFromBridge(t, ethRPCClient, servers, depositAmount, responses)
}
