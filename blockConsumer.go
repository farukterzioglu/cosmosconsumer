package main

import (
	"context"
	"fmt"
	"sync"
	"time"

	"github.com/syndtr/goleveldb/leveldb"

	rpchttp "github.com/tendermint/tendermint/rpc/client/http"
	"github.com/tendermint/tendermint/types"
)

// BlockConsumer is a consumer task for blocks
type BlockConsumer struct {
	Db *leveldb.DB
}

// NewBlockConsumer method creates a new BlockConsumer
func NewBlockConsumer(db *leveldb.DB) *BlockConsumer {
	return &BlockConsumer{
		Db: db,
	}
}

// ConsumeBlocks starts
func (consumer *BlockConsumer) ConsumeBlocks(ctx context.Context, wg *sync.WaitGroup) {
	defer wg.Done()

	client, err := rpchttp.New("tcp://52.137.47.46:26657", "/websocket")
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
		case e := <-txs:
			txHash := e.Events["tx.hash"][0]
			_ = e.Events["tx.height"][0]
			sender := e.Events["transfer.sender"][0]
			recipient := e.Events["transfer.recipient"][0]
			amount := e.Events["transfer.amount"][0]

			data := e.Data.(types.EventDataTx)
			fmt.Printf("Txhash: %s, Height: %d, Index: %d \n", txHash, data.Height, data.Index)
			fmt.Printf("From: %s, To: %s, Amount: %s \n", sender, recipient, amount)
			fmt.Printf("Code: %d, Gas used: %d \n", data.Result.Code, data.Result.GasUsed)

			if false {
				fmt.Println("Events: ")
				for i, oneEvent := range data.Result.Events {
					fmt.Printf("event[%d]: type: %s \n", i, oneEvent.Type)
					for _, attr := range oneEvent.Attributes {
						fmt.Printf("Attribute: %s : %s \n", attr.Key, attr.Value)
					}
				}
			}
		case currenctTime := <-time.After(10 * time.Minute):
			fmt.Printf("Waiting for new blocks (%s)... \n", currenctTime)
		case <-ctx.Done():
			fmt.Println("Block consuming stopped.")
			return
		}
	}
}
