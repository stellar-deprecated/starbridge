package client

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
	"strings"
	"time"

	"github.com/ethereum/go-ethereum"

	"github.com/stellar/starbridge/controllers"

	"github.com/stellar/go/support/render/problem"

	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/crypto"
	"github.com/ethereum/go-ethereum/ethclient"

	"github.com/stellar/go/clients/horizonclient"
	"github.com/stellar/go/keypair"
	"github.com/stellar/go/protocols/horizon"
	"github.com/stellar/go/strkey"
	"github.com/stellar/go/txnbuild"
	"github.com/stellar/starbridge/solidity-go"
)

type BridgeClient struct {
	ValidatorURLs               []string
	EthereumURL                 string
	EthereumChainID             int
	HorizonURL                  string
	NetworkPassphrase           string
	EthereumBridgeAddress       string
	StellarBridgeAccount        string
	EthereumBridgeConfigVersion uint32
	StellarPrivateKey           string
	EthereumPrivateKey          string
}

func (b BridgeClient) SubmitStellarDeposit(amount, asset, ethereumRecipient string) (*horizon.Transaction, error) {
	clientKey := keypair.MustParseFull(b.StellarPrivateKey)
	horizonClient := &horizonclient.Client{
		HorizonURL: b.HorizonURL,
	}

	account, err := horizonClient.AccountDetail(horizonclient.AccountRequest{
		AccountID: clientKey.Address(),
	})
	if err != nil {
		return nil, err
	}

	var txAsset txnbuild.Asset
	if asset == "native" {
		txAsset = txnbuild.NativeAsset{}
	} else {
		parts := strings.Split(asset, ":")
		txAsset = txnbuild.CreditAsset{
			Code:   parts[0],
			Issuer: parts[1],
		}
	}

	tx, err := txnbuild.NewTransaction(txnbuild.TransactionParams{
		SourceAccount: &account,
		Operations: []txnbuild.Operation{
			&txnbuild.Payment{
				Destination: b.StellarBridgeAccount,
				Amount:      amount,
				Asset:       txAsset,
			},
		},
		BaseFee: txnbuild.MinBaseFee,
		Memo:    txnbuild.MemoHash(common.HexToAddress(ethereumRecipient).Hash()),
		Preconditions: txnbuild.Preconditions{
			TimeBounds: txnbuild.NewInfiniteTimeout(),
		},
		IncrementSequenceNum: true,
	})
	if err != nil {
		return nil, err
	}

	tx, err = tx.Sign(b.NetworkPassphrase, clientKey)
	if err != nil {
		return nil, err
	}

	return b.submitStellarTx(horizonClient, tx)
}

func (b BridgeClient) SubmitEthereumDeposit(
	ctx context.Context,
	token common.Address,
	stellarRecipient string,
	amount,
	gasPrice *big.Int,
) (*types.Receipt, error) {
	ethRPCClient, bridge, opts, err := b.createEthClient(gasPrice)
	if err != nil {
		return nil, err
	}

	rawRecipient := strkey.MustDecode(strkey.VersionByteAccountID, stellarRecipient)
	recipient := &big.Int{}
	recipient.SetBytes(rawRecipient)

	var tx *types.Transaction
	if token == (common.Address{}) {
		opts.Value = amount
		tx, err = bridge.DepositETH(opts, recipient)
	} else {
		tx, err = bridge.DepositERC20(opts, token, recipient, amount)
	}
	if err != nil {
		return nil, err
	}

	return submitEthereumTx(ctx, ethRPCClient, tx)
}

func (b BridgeClient) createEthClient(gasPrice *big.Int) (*ethclient.Client, *solidity.BridgeTransactor, *bind.TransactOpts, error) {
	parsedPrivateKey, err := crypto.HexToECDSA(b.EthereumPrivateKey)
	if err != nil {
		return nil, nil, nil, err
	}

	ethRPCClient, err := ethclient.Dial(b.EthereumURL)
	if err != nil {
		return nil, nil, nil, err
	}

	bridge, err := solidity.NewBridgeTransactor(common.HexToAddress(b.EthereumBridgeAddress), ethRPCClient)
	if err != nil {
		return nil, nil, nil, err
	}

	opts, err := bind.NewKeyedTransactorWithChainID(parsedPrivateKey, big.NewInt(int64(b.EthereumChainID)))
	if err != nil {
		return nil, nil, nil, err
	}
	opts.GasPrice = gasPrice
	return ethRPCClient, bridge, opts, nil
}

func (b BridgeClient) SubmitEthereumWithdrawal(
	ctx context.Context,
	stellarTxHash string,
	gasPrice *big.Int,
) (*types.Receipt, error) {
	postData := url.Values{
		"transaction_hash": {stellarTxHash},
	}
	return b.withdrawEthereum(ctx, "stellar/withdraw", postData, gasPrice)
}

func (b BridgeClient) SubmitEthereumRefund(
	ctx context.Context,
	ethereumTxHash string,
	logIndex uint,
	gasPrice *big.Int,
) (*types.Receipt, error) {
	postData := url.Values{
		"transaction_hash": {ethereumTxHash},
		"log_index":        {strconv.FormatUint(uint64(logIndex), 10)},
	}
	return b.withdrawEthereum(ctx, "ethereum/refund", postData, gasPrice)
}

func (b BridgeClient) withdrawEthereum(
	ctx context.Context,
	uri string,
	postData url.Values,
	gasPrice *big.Int,
) (*types.Receipt, error) {
	responses, err := b.ethereumSignatures(uri, postData)
	if err != nil {
		return nil, err
	}

	ethRPCClient, bridge, opts, err := b.createEthClient(gasPrice)
	if err != nil {
		return nil, err
	}

	caller, err := solidity.NewBridgeCaller(common.HexToAddress(b.EthereumBridgeAddress), ethRPCClient)
	if err != nil {
		return nil, err
	}

	validatorToIndex := map[common.Address]uint8{}
	for i := 0; i < len(b.ValidatorURLs); i++ {
		address, err := caller.Signers(nil, big.NewInt(int64(i)))
		if err != nil {
			return nil, err
		}
		validatorToIndex[address] = uint8(i)
	}

	sort.Slice(responses, func(i, j int) bool {
		address := common.HexToAddress(responses[i].Address)
		index := validatorToIndex[address]

		otherAddress := common.HexToAddress(responses[j].Address)
		otherIndex := validatorToIndex[otherAddress]

		return index < otherIndex
	})
	signatures := make([][]byte, len(responses))
	indexes := make([]uint8, len(responses))
	for i, response := range responses {
		signatures[i] = common.Hex2Bytes(response.Signature)
		indexes[i] = validatorToIndex[common.HexToAddress(response.Address)]
	}

	token := common.HexToAddress(responses[0].Token)
	amount, ok := new(big.Int).SetString(responses[0].Amount, 10)
	if !ok {
		return nil, fmt.Errorf("invalid amount in response %v", responses[0].Amount)
	}
	var tx *types.Transaction
	if token == (common.Address{}) {
		tx, err = bridge.WithdrawETH(
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
	} else {
		tx, err = bridge.WithdrawERC20(
			opts,
			solidity.WithdrawERC20Request{
				Id:         common.HexToHash(responses[0].DepositID),
				Expiration: big.NewInt(responses[0].Expiration),
				Recipient:  opts.From,
				Amount:     amount,
				Token:      token,
			},
			signatures,
			indexes,
		)
	}
	if err != nil {
		return nil, err
	}

	return submitEthereumTx(ctx, ethRPCClient, tx)
}

func (b BridgeClient) ethereumSignatures(uri string, postData url.Values) ([]controllers.EthereumSignatureResponse, error) {
	responses := make([]controllers.EthereumSignatureResponse, len(b.ValidatorURLs))
	for i := 0; i < len(b.ValidatorURLs); i++ {
		for {
			requestURL := strings.TrimSuffix(b.ValidatorURLs[i], "/") + "/" + strings.TrimPrefix(uri, "/")
			resp, err := http.PostForm(requestURL, postData)
			if err != nil {
				return nil, err
			}
			switch resp.StatusCode {
			case http.StatusAccepted:
				time.Sleep(time.Second)
				continue
			case http.StatusOK:
				if err = json.NewDecoder(resp.Body).Decode(&responses[i]); err != nil {
					return nil, err
				}
			default:
				return nil, b.parseProblem(resp)
			}
			break
		}
	}
	return responses, nil
}

func (b BridgeClient) parseProblem(resp *http.Response) error {
	var p problem.P
	if err := json.NewDecoder(resp.Body).Decode(&p); err != nil {
		return err
	}
	return p
}

func (b BridgeClient) SubmitStellarRefund(
	stellarTxHash string,
) (*horizon.Transaction, error) {
	postData := url.Values{
		"transaction_hash": {stellarTxHash},
	}
	tx, err := b.stellarTx("stellar/refund", postData)
	if err != nil {
		return nil, err
	}

	horizonClient := &horizonclient.Client{
		HorizonURL: b.HorizonURL,
	}
	return b.submitStellarTx(horizonClient, tx)
}

func (b BridgeClient) SubmitStellarWithdrawal(
	ethereumTxHash string,
	logIndex uint,
) (*horizon.Transaction, error) {
	postData := url.Values{
		"transaction_hash": {ethereumTxHash},
		"log_index":        {strconv.FormatUint(uint64(logIndex), 10)},
	}
	tx, err := b.stellarTx("ethereum/withdraw", postData)
	if err != nil {
		return nil, err
	}

	horizonClient := &horizonclient.Client{
		HorizonURL: b.HorizonURL,
	}
	return b.submitStellarTx(horizonClient, tx)
}

func (b BridgeClient) submitStellarTx(horizonClient *horizonclient.Client, tx *txnbuild.Transaction) (*horizon.Transaction, error) {
	result, err := horizonClient.SubmitTransaction(tx)
	if err != nil {
		return nil, err
	}
	if !result.Successful {
		return nil, fmt.Errorf("transaction %v not successful", result.Hash)
	}
	return &result, nil
}

func (b BridgeClient) stellarTx(uri string, postData url.Values) (*txnbuild.Transaction, error) {
	responses := make([]string, len(b.ValidatorURLs))
	for i := 0; i < len(b.ValidatorURLs); i++ {
		for {
			requestURL := strings.TrimSuffix(b.ValidatorURLs[i], "/") + "/" + strings.TrimPrefix(uri, "/")
			resp, err := http.PostForm(requestURL, postData)
			if err != nil {
				return nil, err
			}
			switch resp.StatusCode {
			case http.StatusAccepted:
				time.Sleep(time.Second)
				continue
			case http.StatusOK:
				if body, err := ioutil.ReadAll(resp.Body); err != nil {
					return nil, err
				} else {
					responses[i] = string(body)
				}
			default:
				return nil, b.parseProblem(resp)
			}
			break
		}
	}

	gtx, err := txnbuild.TransactionFromXDR(responses[0])
	if err != nil {
		return nil, err
	}

	mainTx, ok := gtx.Transaction()
	if !ok {
		return nil, fmt.Errorf("invalid transaction type")
	}

	// Add as many sigs as needed and not a single more
	for i := 1; i < len(b.ValidatorURLs)/2+1; i++ {
		gtx, err := txnbuild.TransactionFromXDR(responses[i])
		if err != nil {
			return nil, err
		}
		tx, ok := gtx.Transaction()
		if !ok {
			return nil, fmt.Errorf("invalid transaction type")
		}

		sig := tx.Signatures()

		mainTx, err = mainTx.AddSignatureDecorated(sig...)
		if err != nil {
			return nil, err
		}
	}

	// ...and add client signature (because it's tx source)
	clientKey := keypair.MustParseFull(b.StellarPrivateKey)
	return mainTx.Sign(b.NetworkPassphrase, clientKey)
}

func submitEthereumTx(ctx context.Context, ethRPCClient *ethclient.Client, tx *types.Transaction) (*types.Receipt, error) {
	var receipt *types.Receipt
	var err error
	sleepDuration := 5 * time.Second

	for attempts := 0; attempts < int(time.Minute/sleepDuration); attempts++ {
		receipt, err = ethRPCClient.TransactionReceipt(ctx, tx.Hash())
		if err == ethereum.NotFound {
			time.Sleep(sleepDuration)
			continue
		} else if err != nil {
			return nil, err
		}
		break
	}
	if err != nil {
		return receipt, err
	}
	if receipt.Status != 1 {
		return nil, fmt.Errorf("transaction %v failed", tx.Hash().String())
	}
	return receipt, nil
}
