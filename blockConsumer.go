package main

import (
	"context"
	"encoding/json"
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"

	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	ctypes "github.com/tendermint/tendermint/rpc/core/types"
	"github.com/tendermint/tendermint/types"
)

// BlockConsumer is a consumer task for blocks
type BlockConsumer struct {
	Db     *leveldb.DB
	Server string
}

// NewBlockConsumer method creates a new BlockConsumer
func NewBlockConsumer(db *leveldb.DB, server string) *BlockConsumer {
	return &BlockConsumer{
		Db:     db,
		Server: server,
	}
}

// ConsumeBlocks starts
func (consumer *BlockConsumer) ConsumeBlocks(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	client, err := rpchttp.New("tcp://"+*server, "/websocket")
	if err != nil {
		fmt.Printf("Got error while creating client, exception: %s \n", err.Error())
		return
	}

	err = client.Start()
	if err != nil {
		fmt.Printf("Got error while starting websocket server, exception: %s \n", err.Error())
		return
	}
	defer client.Stop()

	query := "tm.event='Tx' And message.module='bank' And message.action ='send' And transfer.sender='cosmos17yj45mrgvwezj9jlnhcgd2wr33ufqlcflj0xxv'"
	txs, err := client.Subscribe(ctx, "consumer", query)
	if err != nil {
		fmt.Printf("Got error while subscribing to events , exception: %s \n", err.Error())
		return
	}

	for {
		select {
		case tx := <-txs:
			processTransaction(tx)
		case currenctTime := <-time.After(10 * time.Minute):
			fmt.Printf("Waiting for new blocks (%s)... \n", currenctTime)
		case <-ctx.Done():
			fmt.Println("Block consuming stopped.")
			return
		}
	}
}

func processTransaction(tx ctypes.ResultEvent) {
	// TODO: validate tx (e.g. is success)
	txHash := tx.Events["tx.hash"][0]
	height, _ := strconv.Atoi(tx.Events["tx.height"][0])
	sender := tx.Events["transfer.sender"][0]
	recipient := tx.Events["transfer.recipient"][0]
	amount := tx.Events["transfer.amount"][0]

	data := tx.Data.(types.EventDataTx)
	gasUsed := data.Result.GasUsed

	if false {
		fmt.Printf("Txhash: %s, Height: %d, Index: %d \n", txHash, data.Height, data.Index)
		fmt.Printf("From: %s, To: %s, Amount: %s \n", sender, recipient, amount)
		fmt.Printf("Code: %d, Gas used: %d \n", data.Result.Code, gasUsed)
	}

	if false {
		fmt.Println("Events: ")
		for i, oneEvent := range data.Result.Events {
			fmt.Printf("event[%d]: type: %s \n", i, oneEvent.Type)
			for _, attr := range oneEvent.Attributes {
				fmt.Printf("Attribute: %s : %s \n", attr.Key, attr.Value)
			}
		}
	}

	newTx := Transaction{
		TxHash:    txHash,
		Height:    height,
		Sender:    sender,
		Recipient: recipient,
		Amount:    amount,
		GasUsed:   gasUsed,
	}
	serializedTx, _ := json.Marshal(newTx)
	fmt.Printf(string(serializedTx))
}

// Transaction represents the model to store on leveldb
type Transaction struct {
	TxHash    string `json:"tx_hash,omitempty"`
	Height    int    `json:"height,omitempty"`
	Sender    string `json:"sender,omitempty"`
	Recipient string `json:"recipient,omitempty"`
	Amount    string `json:"amount,omitempty"`
	GasUsed   int64  `json:"gas_used,omitempty"`
}
