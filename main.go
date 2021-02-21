package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"

	"github.com/syndtr/goleveldb/leveldb"
)

var (
	dataDir    = flag.String("data-dir", ".gaiaconsumer", "")
	server     = flag.String("server", "0.0.0.0:26657", "host:port")
	startHeigh = flag.String("start-height", "-1", "")
)

const (
	keyStartingHeight string = "startingHeight"
)

const (
	blockPrefix       string = "block_"
	transactionPrefix string = "transaction_"
)

func main() {
	flag.Parse()

	fmt.Printf("Data directory: %s\n", *dataDir)
	db, err := leveldb.OpenFile(fmt.Sprintf("%s/db", *dataDir), nil)
	if err != nil {
		fmt.Printf("Error while openning leveldb: %s", err.Error())
		return
	}
	defer db.Close()

	startingHeight, _ := strconv.Atoi(*startHeigh)
	if startingHeight == -1 {
		keyStartingHeightByte := []byte(keyStartingHeight)

		hasKeyStartingHeight, err := db.Has(keyStartingHeightByte, nil)
		if err != nil {
			fmt.Printf("Couldn't get value for %s, exception: %s.\n", keyStartingHeight, err.Error())
			return
		}

		if !hasKeyStartingHeight {
			err = db.Put(keyStartingHeightByte, []byte("0"), nil)
			if err != nil {
				fmt.Printf("Couldn't put value for %s, exception: %s.\n", keyStartingHeight, err.Error())
				return
			}
			fmt.Printf("Setting starting height as 0 in db.\n")
		}

		data, err := db.Get(keyStartingHeightByte, nil)
		if err != nil {
			fmt.Printf("Couldn't get value for %s, exception: %s.\n", keyStartingHeight, err.Error())
			return
		}

		startingHeight, _ = strconv.Atoi(string(data))
	}
	fmt.Printf("Starting from block: %d\n", startingHeight)

	ctx, cancelFunc := context.WithCancel(context.Background())
	wg := &sync.WaitGroup{}

	fmt.Printf("Connecting to server: %s\n", *server)
	var consumer *BlockConsumer
	consumer = NewBlockConsumer(db, *server)
	wg.Add(1)
	go consumer.ConsumeBlocks(ctx, wg)

	termChan := make(chan os.Signal)
	signal.Notify(termChan, syscall.SIGINT, syscall.SIGTERM)
	<-termChan

	cancelFunc()
	wg.Wait()

	fmt.Println("Consumer stopped.")
}
