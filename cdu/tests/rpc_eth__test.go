package tests

import (
	"fmt"
	"github.com/RevittConsulting/chain-dev-utils/internal/jsonrpc"
	"github.com/magiconair/properties/assert"
	"math/big"
	"testing"
)

func Test_EthBlockNumber(t *testing.T) {
	client := jsonrpc.NewClient("https://rpc.sepolia.org")

	blockNum, err := client.EthBlockNumber()
	if err != nil {
		t.Error(err)
	}

	t.Log("blockNum:", blockNum)

	assert.Equal(t, blockNum > 0, true)
}

func Test_EthGetLogs(t *testing.T) {
	client := jsonrpc.NewClient("https://rpc.sepolia.org")

	blockNum, err := client.EthBlockNumber()
	if err != nil {
		t.Error(err)
	}

	fromBlock := fmt.Sprintf("0x%X", blockNum-2000)
	toBlock := "latest"
	address := interface{}("0x32d33D5137a7cFFb54c5Bf8371172bcEc5f310ff")

	query := jsonrpc.LogQuery{
		FromBlock: &fromBlock,
		ToBlock:   &toBlock,
		Address:   &address,
		Topics:    nil,
	}

	logs, err := client.EthGetLogs(query)
	if err != nil {
		t.Error(err)
	}

	t.Log("logs:", logs)

	assert.Equal(t, len(logs) > 0, true)
}

func Test_EthGasPrice(t *testing.T) {
	client := jsonrpc.NewClient("https://rpc.sepolia.org")

	gasPrice, err := client.EthGasPrice()
	if err != nil {
		t.Error(err)
	}

	t.Log("gasPrice:", gasPrice)

	assert.Equal(t, gasPrice.Cmp(new(big.Int).SetUint64(0)) > 0, true)
}
