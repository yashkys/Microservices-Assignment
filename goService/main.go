package main

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"os/signal"
	"time"

	"encoding/json"

	"ass1.com/transaction"
	"github.com/gorilla/websocket"
	"github.com/joho/godotenv"
	"github.com/ybbus/jsonrpc/v3"

	"google.golang.org/grpc"
)

type BlockResult struct {
	BaseFeePerGas    string   `json:"baseFeePerGas"`
	Difficulty       string   `json:"difficulty"`
	ExtraData        string   `json:"extraData"`
	GasLimit         string   `json:"gasLimit"`
	GasUsed          string   `json:"gasUsed"`
	Hash             string   `json:"hash"`
	L1BlockNumber    string   `json:"l1BlockNumber"`
	LogsBloom        string   `json:"logsBloom"`
	Miner            string   `json:"miner"`
	MixHash          string   `json:"mixHash"`
	Nonce            string   `json:"nonce"`
	Number           string   `json:"number"`
	ParentHash       string   `json:"parentHash"`
	ReceiptsRoot     string   `json:"receiptsRoot"`
	SendCount        string   `json:"sendCount"`
	SendRoot         string   `json:"sendRoot"`
	Sha3Uncles       string   `json:"sha3Uncles"`
	Size             string   `json:"size"`
	StateRoot        string   `json:"stateRoot"`
	Timestamp        string   `json:"timestamp"`
	TotalDifficulty  string   `json:"totalDifficulty"`
	Transactions     []string `json:"transactions"`
	TransactionsRoot string   `json:"transactionsRoot"`
	Uncles           []string `json:"uncles"`
}
type Log struct {
	BlockHash        string   `json:"blockHash"`
	Address          string   `json:"address"`
	LogIndex         string   `json:"logIndex"`
	Data             string   `json:"data"`
	Removed          bool     `json:"removed"`
	Topics           []string `json:"topics"`
	BlockNumber      string   `json:"blockNumber"`
	TransactionIndex string   `json:"transactionIndex"`
	TransactionHash  string   `json:"transactionHash"`
}
type TransactionReceipt struct {
	TransactionHash   string `json:"transactionHash"`
	BlockHash         string `json:"blockHash"`
	BlockNumber       string `json:"blockNumber"`
	LogsBloom         string `json:"logsBloom"`
	L1BlockNumber     string `json:"l1BlockNumber"`
	ContractAddress   string `json:"contractAddress"`
	TransactionIndex  string `json:"transactionIndex"`
	Type              string `json:"type"`
	GasUsed           string `json:"gasUsed"`
	CumulativeGasUsed string `json:"cumulativeGasUsed"`
	From              string `json:"from"`
	To                string `json:"to"`
	EffectiveGasPrice string `json:"effectiveGasPrice"`
	Logs              []Log  `json:"logs"`
	Status            string `json:"status"`
	GasUsedForL1      string `json:"gasUsedForL1"`
}

type BlockHeader struct {
	Difficulty       string `json:"difficulty"`
	ExtraData        string `json:"extraData"`
	GasLimit         string `json:"gasLimit"`
	GasUsed          string `json:"gasUsed"`
	LogsBloom        string `json:"logsBloom"`
	Miner            string `json:"miner"`
	Nonce            string `json:"nonce"`
	Number           string `json:"number"`
	ParentHash       string `json:"parentHash"`
	ReceiptsRoot     string `json:"receiptRoot"`
	Sha3Uncles       string `json:"sha3Uncles"`
	StateRoot        string `json:"stateRoot"`
	Timestamp        string `json:"timestamp"`
	TransactionsRoot string `json:"transactionsRoot"`
}

type SubscriptionParams struct {
	Result       BlockHeader `json:"result"`
	Subscription string      `json:"subscription"`
}

type SubscriptionMessage struct {
	Jsonrpc string             `json:"jsonrpc"`
	Method  string             `json:"method"`
	Params  SubscriptionParams `json:"params"`
}

func main() {
	//initalize logger
	initLogger()
	err := godotenv.Load()
	if err != nil {
		errorLog.Println("Error loading .env file")
	}

	// Channels for messaging and interruptions
	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)
	apiKey := os.Getenv("API_KEY")

	// WebSocket URL setup
	u := url.URL{
		Scheme: "wss",
		Host:   "arb-mainnet.g.alchemy.com",
		Path:   fmt.Sprintf("/v2/%s", apiKey),
	}

	// Establish WebSocket connection
	c, resp, err := websocket.DefaultDialer.Dial(u.String(), nil)
	if err != nil {
		errorLog.Printf("Failed to establish websocket connection(Status : %d, Error : %s)\n", resp.StatusCode, err)
	}
	successLog.Println("Websocket connection established with ", u.String())
	defer c.Close()

	// Handle incoming messages in a separate goroutine
	done := make(chan struct{})
	go func() {
		defer close(done)
		for {
			_, message, err := c.ReadMessage()
			if err != nil {
				errorLog.Println("Failed to read message: ", err)
				return
			}

			// Handle incoming block header updates
			var subMessage SubscriptionMessage
			err = json.Unmarshal(message, &subMessage)
			if err != nil {
				errorLog.Println("Invalid Json format received:", err)
				return
			}

			if subMessage.Method == "eth_subscription" {
				blockNumberHex := subMessage.Params.Result.Number
				infoLog.Println("Block number received: ", blockNumberHex)
				performOperationOnBlock(blockNumberHex, apiKey)
			}
		}
	}()

	// Subscribe to new block headers
	subscribeMessage := `{"jsonrpc":"2.0","id":1,"method":"eth_subscribe","params":["newHeads"]}`
	err = c.WriteMessage(websocket.TextMessage, []byte(subscribeMessage))
	if err != nil {
		errorLog.Println("Failed to write message:", err)
		return
	}

	// Set up a ticker to keep the connection alive
	ticker := time.NewTicker(time.Second * 10)
	defer ticker.Stop()

	for {
		select {
		case <-done:
			return
		case t := <-ticker.C:
			// Send ping to keep the connection alive
			err := c.WriteMessage(websocket.TextMessage, []byte(t.String()))
			if err != nil {
				errorLog.Println("Failed to write message:", err)
				return
			}
		case <-interrupt:
			warningLog.Println("Interrupt")
			// Cleanly close the connection
			err := c.WriteMessage(websocket.CloseMessage, websocket.FormatCloseMessage(websocket.CloseNormalClosure, ""))
			if err != nil {
				errorLog.Println("Failed to close write connection:", err)
				return
			}
			select {
			case <-done:
			case <-time.After(time.Second):
			}
			return
		}
	}
}

func performOperationOnBlock(blockNumberHex string, apiKey string) {
	// Create a JSON-RPC client
	alchemyUrl := fmt.Sprintf("https://arb-mainnet.g.alchemy.com/v2/%s", apiKey)
	rpcClient := jsonrpc.NewClient(alchemyUrl)
	service2url := "localhost:50051"
	grpcConnection, err := grpc.NewClient(service2url, grpc.WithInsecure())
	if err != nil {
		errorLog.Printf("GRPC connection error with %s: %s\n", service2url, err)
	}
	defer grpcConnection.Close()
	successLog.Println("Grpc connection established with ", service2url)
	client := transaction.NewTransactionServiceClient(grpcConnection)

	// Make a POST call to retrieve the latest block
	response, err := rpcClient.Call(context.Background(), "eth_getBlockByNumber", blockNumberHex, false)
	if err != nil {
		errorLog.Println("Failed to fetch block: ", err)
	}

	// Unmarshal the response into the BlockResult struct
	var blockResult BlockResult
	if err := response.GetObject(&blockResult); err != nil {
		errorLog.Println("Failed to parse block: ", err)
	}

	// Iterate over each transaction hash and get the transaction receipt
	for _, txHash := range blockResult.Transactions {
		// Make the JSON-RPC call to get the transaction receipt
		receiptResponse, err := rpcClient.Call(context.Background(), "eth_getTransactionReceipt", txHash)
		if err != nil {
			errorLog.Printf("Failed to fetch transaction receipt for txHash %s: %v \n", txHash, err)
			continue
		}
		var receipt TransactionReceipt
		err = receiptResponse.GetObject(&receipt)
		if err != nil {
			errorLog.Printf("Failed to unmarshal transaction receipt for txHash %s: %v\n", txHash, err)
			continue
		}

		var receiptProto transaction.TransactionReceipt

		receiptProto.TransactionHash = txHash
		receiptProto.BlockNumber = receipt.BlockNumber
		receiptProto.BlockHash = receipt.BlockHash
		receiptProto.LogsBloom = receipt.LogsBloom
		receiptProto.L1BlockNumber = receipt.L1BlockNumber
		receiptProto.ContractAddress = receipt.ContractAddress
		receiptProto.TransactionIndex = receipt.TransactionIndex
		receiptProto.Type = receipt.Type
		receiptProto.GasUsed = receipt.GasUsed
		receiptProto.CumulativeGasUsed = receipt.CumulativeGasUsed
		receiptProto.From = receipt.From
		receiptProto.To = receipt.To
		receiptProto.EffectiveGasPrice = receipt.EffectiveGasPrice
		receiptProto.Logs = mapperToMapLogsInTransactionReceipt(receipt.Logs)
		receiptProto.Status = receipt.Status
		receiptProto.GasUsedForL1 = receipt.GasUsedForL1

		_, err = client.SubmitTransactionReceipt(context.Background(), &receiptProto)
		if err != nil {
			errorLog.Println("Internal Server Error: ", err)
		}
		// successLog.Println("TxHash transaction receipt sent")
	}
}

func mapperToMapLogsInTransactionReceipt(logs []Log) []*transaction.Log {
	var protoLogs []*transaction.Log
	for _, log := range logs {
		protoLogs = append(protoLogs, &transaction.Log{
			BlockHash:        log.BlockHash,
			Address:          log.Address,
			LogIndex:         log.LogIndex,
			Data:             log.Data,
			Removed:          log.Removed,
			Topics:           log.Topics,
			BlockNumber:      log.BlockNumber,
			TransactionIndex: log.TransactionIndex,
			TransactionHash:  log.TransactionHash,
		})
	}
	return protoLogs
}
