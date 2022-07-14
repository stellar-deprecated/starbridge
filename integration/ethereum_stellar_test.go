package integration

import (
	"context"
	"database/sql"
	"encoding/base64"
	"encoding/hex"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"

	"github.com/stretchr/testify/require"

	"github.com/stellar/starbridge/store"
)

var gasPrice = new(big.Int).Mul(big.NewInt(1), big.NewInt(params.GWei))

func ethereumSenderAddress(t *testing.T) common.Address {
	parsedPrivateKey, err := crypto.HexToECDSA(ethereumSenderPrivateKey)
	require.NoError(t, err)
	ethereumRecipient := crypto.PubkeyToAddress(parsedPrivateKey.PublicKey)
	return ethereumRecipient
}

func TestEthereumToStellarWithdrawal(t *testing.T) {
	const servers int = 3

	itest := NewIntegrationTest(t, Config{
		Servers:                servers,
		EthereumFinalityBuffer: 0,
		WithdrawalWindow:       time.Hour,
	})

	receipt, err := itest.bridgeClient.SubmitEthereumDeposit(
		context.Background(),
		common.Address{},
		itest.clientKey.Address(),
		new(big.Int).Mul(big.NewInt(3), big.NewInt(params.Ether)),
		gasPrice,
	)
	require.NoError(t, err)

	_, err = itest.bridgeClient.SubmitEthereumRefund(
		context.Background(),
		receipt.TxHash.String(),
		receipt.Logs[0].Index,
		nil,
	)
	require.EqualError(t, err, "problem: https://stellar.org/horizon-errors/withdrawal_window_still_active")

	tx, err := itest.bridgeClient.SubmitStellarWithdrawal(receipt.TxHash.String(), receipt.Logs[0].Index)
	require.NoError(t, err)
	memoBytes, err := base64.StdEncoding.DecodeString(tx.Memo)
	require.NoError(t, err)

	stores := make([]*store.DB, servers)
	for i := 0; i < servers; i++ {
		stores[i] = itest.app[i].NewStore()
	}
	numFound := 0
	for attempts := 0; attempts < 10; attempts++ {
		for i := 0; i < servers; i++ {
			found, err := stores[i].HistoryStellarTransactionExists(context.Background(), hex.EncodeToString(memoBytes))
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

	_, err = itest.bridgeClient.SubmitStellarWithdrawal(receipt.TxHash.String(), receipt.Logs[0].Index)
	require.EqualError(t, err, "problem: https://stellar.org/horizon-errors/withdrawal_already_executed")
}

func TestEthereumRefund(t *testing.T) {
	const servers int = 3

	itest := NewIntegrationTest(t, Config{
		Servers:                servers,
		EthereumFinalityBuffer: 0,
		WithdrawalWindow:       time.Second,
	})

	depositAmount := new(big.Int).Mul(big.NewInt(3), big.NewInt(params.Ether))
	receipt, err := itest.bridgeClient.SubmitEthereumDeposit(
		context.Background(),
		common.Address{},
		itest.clientKey.Address(),
		depositAmount,
		gasPrice,
	)
	require.NoError(t, err)

	stores := make([]*store.DB, servers)
	for i := 0; i < servers; i++ {
		stores[i] = itest.app[i].NewStore()
	}

	ethRPCClient, err := ethclient.Dial(EthereumRPCURL)
	require.NoError(t, err)
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

	_, err = itest.bridgeClient.SubmitStellarWithdrawal(receipt.TxHash.String(), receipt.Logs[0].Index)
	require.EqualError(t, err, "problem: https://stellar.org/horizon-errors/withdrawal_window_expired")

	_, err = itest.bridgeClient.SubmitEthereumRefund(
		context.Background(),
		receipt.TxHash.String(),
		receipt.Logs[0].Index,
		gasPrice,
	)
	require.NoError(t, err)
}

func TestStellarToEthereumWithdrawal(t *testing.T) {
	const servers int = 3

	itest := NewIntegrationTest(t, Config{
		Servers:                servers,
		EthereumFinalityBuffer: 0,
		WithdrawalWindow:       time.Hour,
	})

	tx, err := itest.bridgeClient.SubmitStellarDeposit("3", "native", ethereumSenderAddress(t).String())
	require.NoError(t, err)

	stores := make([]*store.DB, servers)
	for i := 0; i < servers; i++ {
		stores[i] = itest.app[i].NewStore()
	}
	for {
		ready := 0
		for i := 0; i < servers; i++ {
			_, err = stores[i].GetStellarDeposit(context.Background(), tx.Hash)
			if err == sql.ErrNoRows {
				continue
			}
			require.NoError(t, err)
			ready++
		}
		if ready == servers {
			break
		}
		time.Sleep(time.Second)
	}

	_, err = itest.bridgeClient.SubmitStellarRefund(tx.Hash)
	require.EqualError(t, err, "problem: https://stellar.org/horizon-errors/withdrawal_window_still_active")

	_, err = itest.bridgeClient.SubmitEthereumWithdrawal(
		context.Background(),
		tx.Hash,
		gasPrice,
	)
	require.NoError(t, err)
}

func TestStellarRefund(t *testing.T) {
	const servers int = 3

	itest := NewIntegrationTest(t, Config{
		Servers:                servers,
		EthereumFinalityBuffer: 0,
		WithdrawalWindow:       time.Second,
	})

	ethRPCClient, err := ethclient.Dial(EthereumRPCURL)
	require.NoError(t, err)

	tx, err := itest.bridgeClient.SubmitStellarDeposit("3", "native", ethereumSenderAddress(t).String())
	require.NoError(t, err)

	stores := make([]*store.DB, servers)
	for i := 0; i < servers; i++ {
		stores[i] = itest.app[i].NewStore()
	}
	var deposit store.StellarDeposit
	for {
		ready := 0
		for i := 0; i < servers; i++ {
			deposit, err = stores[i].GetStellarDeposit(context.Background(), tx.Hash)
			if err == sql.ErrNoRows {
				continue
			}
			require.NoError(t, err)
			ready++
		}
		if ready == servers {
			break
		}
		time.Sleep(time.Second)
	}

	_, err = itest.bridgeClient.SubmitStellarRefund(tx.Hash)
	require.EqualError(t, err, "problem: https://stellar.org/horizon-errors/withdrawal_window_still_active")

	// Wait for WithdrawalWindow to pass in Ethereum
	depositTime := time.Unix(deposit.LedgerTime, 0)
	for {
		header, err := ethRPCClient.HeaderByNumber(context.Background(), nil)
		require.NoError(t, err)
		t.Log("Block ", header.Number, " time ", header.Time)
		if time.Unix(int64(header.Time), 0).After(depositTime.Add(time.Second)) {
			break
		}

		// Close Ethereum block
		ethClient, err := rpc.DialContext(context.Background(), EthereumRPCURL)
		require.NoError(t, err)
		err = ethClient.Call(nil, "evm_mine")
		require.NoError(t, err)

		time.Sleep(time.Second)
	}

	t.Log("Ethereum time reached withdrawal deadline")

	refundTx, err := itest.bridgeClient.SubmitStellarRefund(tx.Hash)
	require.NoError(t, err)
	memoBytes, err := base64.StdEncoding.DecodeString(refundTx.Memo)
	require.NoError(t, err)

	numFound := 0
	for attempts := 0; attempts < 10; attempts++ {
		for i := 0; i < servers; i++ {
			found, err := stores[i].HistoryStellarTransactionExists(context.Background(), hex.EncodeToString(memoBytes))
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

	_, err = itest.bridgeClient.SubmitStellarRefund(tx.Hash)
	require.EqualError(t, err, "problem: https://stellar.org/horizon-errors/refund_already_executed")
}
