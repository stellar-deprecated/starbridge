package concordium

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	concordiumSDK "github.com/Concordium/concordium-go-sdk"
	concordiumproto "github.com/Concordium/concordium-go-sdk/grpc-api"
	"github.com/stellar/go/support/log"
	"io"
	"io/ioutil"
	"net/http"
	"time"

	"github.com/ethereum/go-ethereum/common"
)

type InvokeWithdrawHashParameters struct {
	Amount     *concordiumSDK.Amount  `json:"amount"`
	Expiration uint64                 `json:"expiration"`
	Id         string                 `json:"id"`
	To         *concordiumSDK.Address `json:"to"`
}
type InvokeWithdrawHashResponse struct {
	Hash []uint8 `json:"hash"`
}
type GetDepositParamsResponse struct {
	Amount      string `json:"amount"`
	Destination string `json:"destination"`
	BlockHash   string `json:"blockHash"`
	From        string `json:"from"`
}

type ConsensusStatusResponse struct {
	BestBlock string `json:"bestBlock"`
}

type ConcordiumDeposit struct {
	Amount    string
	From      string
	To        string
	Cost      string
	Hash      string
	Index     int
	BlockHash string
	BlockTime int64
}

type TransactionEvent struct {
	Amount string
	From   struct {
		Address string
	}
	Tag string
	To  struct {
		Address string
	}
}

type TransactionInfo struct {
	Cost      string
	Hash      string
	Index     int
	Sender    string
	BlockHash string
}

type Block struct {
	hash string
}

type Observer struct {
	client                   concordiumproto.P2PClient
	bridgeAddress            common.Address
	concordiumNodeServiceUrl string
}

func NewObserver(client concordiumproto.P2PClient, bridgeAddress string, concordiumNodeServiceUrl string) (Observer, error) {
	//if !common.IsHexAddress(bridgeAddress) {
	//	return Observer{}, fmt.Errorf("%v is not a valid concordium address", bridgeAddress)
	//}
	bridgeAddressParsed := common.HexToAddress(bridgeAddress)

	//caller, err := solidity.NewBridgeCaller(bridgeAddressParsed, client)
	//if err != nil {
	//	return Observer{}, err
	//}
	//filterer, err := solidity.NewBridgeFilterer(bridgeAddressParsed, client)
	//if err != nil {
	//	return Observer{}, err
	//}

	return Observer{
		client: client,
		//filterer:      filterer,
		//caller:        caller,
		bridgeAddress:            bridgeAddressParsed,
		concordiumNodeServiceUrl: concordiumNodeServiceUrl,
	}, nil
}

func (o Observer) GetBestBlock(ctx context.Context) (Block, error) {
	grpcResponse, err := o.client.GetConsensusStatus(ctx, &concordiumproto.Empty{})
	if err != nil {
		log.Error("Error GetConsensusStatus", err)
		return Block{}, err
	}
	consensusStatusResponse := ConsensusStatusResponse{}
	err = json.Unmarshal([]byte(grpcResponse.Value), &consensusStatusResponse)
	if err != nil {
		log.Error("Error json.Unmarshal", err)
		return Block{}, err
	}

	return Block{consensusStatusResponse.BestBlock}, nil
}

func (o Observer) GetTransactionStatus(ctx context.Context, txHash string) (map[string]interface{}, error) {
	var body map[string]interface{}
	grpcResponse, err := o.client.GetTransactionStatus(ctx, &concordiumproto.TransactionHash{TransactionHash: txHash})
	if err != nil {
		return body, err
	}
	if err = json.Unmarshal([]byte(grpcResponse.Value), &body); err != nil {
		return body, err
	}
	return body, nil
}

func (o Observer) GetDeposit(ctx context.Context, txHash string) (ConcordiumDeposit, error) {
	txBody, err := o.GetTransactionStatus(ctx, txHash)
	if err != nil {
		return ConcordiumDeposit{}, err
	}
	txBodyStatus := txBody["status"].(string)
	for txBodyStatus != "finalized" {
		txBody, err = o.GetTransactionStatus(context.Background(), txHash)
		if err != nil {
			return ConcordiumDeposit{}, err
		}
		txBodyStatus = txBody["status"].(string)
	}

	getDepositParamsRequest := struct {
		Hash string `json:"hash"`
	}{
		Hash: txHash,
	}
	getDepositParamsRequestJson, err := json.Marshal(getDepositParamsRequest)
	if err != nil {
		return ConcordiumDeposit{}, err
	}

	jsonBody := []byte(string(getDepositParamsRequestJson))

	req, err := http.NewRequest("POST", fmt.Sprintf("%s/invokeContract/getDepositParams", o.concordiumNodeServiceUrl), bytes.NewBuffer(jsonBody))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return ConcordiumDeposit{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	response := GetDepositParamsResponse{}
	err = json.Unmarshal(body, &response)

	deposit := ConcordiumDeposit{}
	deposit.Hash = txHash
	deposit.BlockHash = response.BlockHash
	deposit.Amount = response.Amount
	deposit.From = response.From

	blockInfo, err := o.GetBlockInfo(ctx, deposit.BlockHash)
	if err != nil {
		return ConcordiumDeposit{}, err
	}
	var bodyBlock map[string]interface{}
	if err = json.Unmarshal([]byte(blockInfo.Value), &bodyBlock); err != nil {
		return ConcordiumDeposit{}, err
	}
	layout := "2006-01-02T15:04:05Z0700"
	blockTime, err := time.Parse(layout, bodyBlock["blockSlotTime"].(string))
	if err = json.Unmarshal([]byte(blockInfo.Value), &bodyBlock); err != nil {
		return ConcordiumDeposit{}, err
	}
	deposit.BlockTime = blockTime.Unix()

	return deposit, nil
}

func (o Observer) GetBlockInfo(ctx context.Context, blockHash string) (*concordiumproto.JsonResponse, error) {
	grpcResponse, err := o.client.GetBlockInfo(ctx, &concordiumproto.BlockHash{BlockHash: blockHash})
	if err != nil {
		return &concordiumproto.JsonResponse{}, err
	}
	return grpcResponse, nil
}

func (o Observer) InvokeWithdrawHash(amount int, expiration uint64, depositId [32]byte, accountDestination string) ([]uint8, error) {
	invokeWithdrawHashRequest := struct {
		Method     string `json:"method"`
		Parameters struct {
			Amount     int      `json:"amount"`
			Expiration uint64   `json:"expiration"`
			Id         [32]byte `json:"id"`
			Indexes    []int    `json:"indexes"`
			Signatures []int    `json:"signatures"`
			To         struct {
				Account []string `json:"Account"`
			} `json:"to"`
		} `json:"parameters"`
	}{
		Method: "withdraw_hash",
		Parameters: struct {
			Amount     int      `json:"amount"`
			Expiration uint64   `json:"expiration"`
			Id         [32]byte `json:"id"`
			Indexes    []int    `json:"indexes"`
			Signatures []int    `json:"signatures"`
			To         struct {
				Account []string `json:"Account"`
			} `json:"to"`
		}(struct {
			Amount     int
			Expiration uint64
			Id         [32]byte
			Indexes    []int
			Signatures []int
			To         struct {
				Account []string
			}
		}{
			Amount:     amount,
			Expiration: expiration * 1000,
			Id:         depositId,
			Indexes:    []int{},
			Signatures: []int{},
			To: struct {
				Account []string
			}{
				Account: []string{accountDestination},
			},
		}),
	}
	invokeWithdrawHashRequestJson, err := json.Marshal(invokeWithdrawHashRequest)
	if err != nil {
		return []uint8{}, err
	}

	jsonStr := []byte(string(invokeWithdrawHashRequestJson))
	req, err := http.NewRequest("POST", fmt.Sprintf("%s/invokeContract/getWithdrawHash", o.concordiumNodeServiceUrl), bytes.NewBuffer(jsonStr))
	req.Header.Set("X-Custom-Header", "myvalue")
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return []uint8{}, err
	}
	defer func(Body io.ReadCloser) {
		err := Body.Close()
		if err != nil {
		}
	}(resp.Body)

	body, _ := ioutil.ReadAll(resp.Body)
	response := InvokeWithdrawHashResponse{}
	err = json.Unmarshal(body, &response)
	if err != nil {
		return []uint8{}, err
	}
	return response.Hash, nil
}
