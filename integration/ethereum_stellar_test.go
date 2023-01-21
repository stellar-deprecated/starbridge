package integration

import (
	"context"
	"math/big"
	"testing"
	"time"

	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ethereum/go-ethereum/params"
	"github.com/ethereum/go-ethereum/rpc"
	"github.com/stretchr/testify/require"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/xdr"
)

var gasPrice = new(big.Int).Mul(big.NewInt(1), big.NewInt(params.GWei))

func ethereumSenderAddress(t *testing.T) common.Address {
	parsedPrivateKey, err := crypto.HexToECDSA(ethereumSenderPrivateKey)
	require.NoError(t, err)
	ethereumRecipient := crypto.PubkeyToAddress(parsedPrivateKey.PublicKey)
	return ethereumRecipient
}

func TestEthereumToStellarWithdrawal(t *testing.T) {
	itest := NewIntegrationTest(t, Config{
		Servers:                3,
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

	_, err = itest.bridgeClient.SubmitStellarWithdrawal(receipt.TxHash.String(), receipt.Logs[0].Index)
	require.NoError(t, err)

	// double withdrawal fails
	_, err = itest.bridgeClient.SubmitStellarWithdrawal(receipt.TxHash.String(), receipt.Logs[0].Index)
	require.Error(t, err)

	tx, err := itest.bridgeClient.SubmitStellarDeposit(
		xdr.MustNewCreditAsset("ETH", itest.bridgeClient.StellarBridgeAccount),
		"1",
		ethereumSenderAddress(t),
	)
	require.NoError(t, err)

	_, err = itest.bridgeClient.SubmitEthereumWithdrawal(
		context.Background(),
		tx.Hash,
		gasPrice,
	)
	require.NoError(t, err)
}

func TestEthereumRefund(t *testing.T) {
	itest := NewIntegrationTest(t, Config{
		Servers:                3,
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

	ethRPCClient, err := ethclient.Dial(EthereumRPCURL)
	require.NoError(t, err)
	header, err := ethRPCClient.HeaderByHash(context.Background(), receipt.BlockHash)
	require.NoError(t, err)
	depositTime := time.Unix(int64(header.Time), 0)
	for {
		ledgers, err := itest.horizonClient.Ledgers(horizonclient.LedgerRequest{
			Order: "desc",
			Limit: 1,
		})
		require.NoError(t, err)
		if ledgers.Embedded.Records[0].ClosedAt.After(depositTime.Add(time.Second)) {
			break
		}

		time.Sleep(time.Second)
	}

	_, err = itest.bridgeClient.SubmitStellarWithdrawal(receipt.TxHash.String(), receipt.Logs[0].Index)
	require.Error(t, err)

	_, err = itest.bridgeClient.SubmitEthereumRefund(
		context.Background(),
		receipt.TxHash.String(),
		receipt.Logs[0].Index,
		gasPrice,
	)
	require.NoError(t, err)

	// double refund fails
	_, err = itest.bridgeClient.SubmitEthereumRefund(
		context.Background(),
		receipt.TxHash.String(),
		receipt.Logs[0].Index,
		gasPrice,
	)
	require.Error(t, err)
}

func TestStellarToEthereumWithdrawal(t *testing.T) {
	itest := NewIntegrationTest(t, Config{
		Servers:                3,
		EthereumFinalityBuffer: 0,
		WithdrawalWindow:       time.Hour,
	})

	tx, err := itest.bridgeClient.SubmitStellarDeposit(xdr.MustNewNativeAsset(), "3", ethereumSenderAddress(t))
	require.NoError(t, err)

	_, err = itest.bridgeClient.SubmitStellarRefund(tx.Hash)
	require.EqualError(t, err, "problem: https://stellar.org/horizon-errors/withdrawal_window_still_active")

	_, err = itest.bridgeClient.SubmitEthereumWithdrawal(
		context.Background(),
		tx.Hash,
		gasPrice,
	)
	require.NoError(t, err)

	// double withdrawal fails
	_, err = itest.bridgeClient.SubmitEthereumWithdrawal(
		context.Background(),
		tx.Hash,
		gasPrice,
	)
	require.Error(t, err)
}

func TestStellarRefund(t *testing.T) {
	itest := NewIntegrationTest(t, Config{
		Servers:                3,
		EthereumFinalityBuffer: 0,
		WithdrawalWindow:       time.Second,
	})

	tx, err := itest.bridgeClient.SubmitStellarDeposit(xdr.MustNewNativeAsset(), "3", ethereumSenderAddress(t))
	require.NoError(t, err)
	require.True(t, tx.Successful)

	_, err = itest.bridgeClient.SubmitStellarRefund(tx.Hash)
	require.EqualError(t, err, "problem: https://stellar.org/horizon-errors/withdrawal_window_still_active")

	waitPastWithdrawalWindow(t, tx)

	_, err = itest.bridgeClient.SubmitEthereumWithdrawal(context.Background(), tx.Hash, gasPrice)
	require.Error(t, err)

	_, err = itest.bridgeClient.SubmitStellarRefund(tx.Hash)
	require.NoError(t, err)

	// double refund fails
	_, err = itest.bridgeClient.SubmitStellarRefund(tx.Hash)
	require.Error(t, err)
}

func waitPastWithdrawalWindow(t *testing.T, tx *horizon.Transaction) {
	ethRPCClient, err := ethclient.Dial(EthereumRPCURL)
	require.NoError(t, err)

	// Wait for WithdrawalWindow to pass in Ethereum
	depositTime := tx.LedgerCloseTime
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
}
