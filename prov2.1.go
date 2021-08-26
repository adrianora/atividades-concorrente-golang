package main

import (
	"fmt"
	"math/rand"
	"time"
)

type Item struct {
	name         string
	initialValue int
}

type Bid struct {
	item      Item
	bidValue  int
	bidFailed bool
}

func itemsStream() chan Item {
	// informa o número de itens que deseja criar
	// está hardcoded para não alterar assinatura da função
	nItems := 10

	items := make(chan Item, nItems)
	for i := 0; i < nItems; i++ {
		var item *Item
		item = new(Item)
		item.name = "name"
		item.initialValue = i
		items <- *item
	}
	close(items)

	return items
}

func bid(item Item) Bid {
	var bided *Bid
	bided = new(Bid)
	bided.item = item
	bided.bidValue = 1
	bided.bidFailed = false

	// bloqueia o processo por tempo aleatório para simular efeito
	// concorrente de uma requisição demorada
	rand.Seed(time.Now().UnixNano())
	sleepTime := rand.Intn(5)
	fmt.Printf("sleeping for %d seconds in item %+v\n", sleepTime, item)
	time.Sleep(time.Duration(sleepTime) * time.Second)

	return *bided
}

func handle(nServers int) chan Bid {
	itemsCh := itemsStream()
	bidsCh := make(chan Bid, len(itemsCh))
	// canal responsável por controlar a existência de todas as goroutines
	awaitServerWorkCh := make(chan bool, nServers)

	for i := 1; i <= nServers; i++ {
		// cria uma goroutine para  cada servidor
		go func() {
			for {
				item, notClosed := <-itemsCh
				if notClosed {
					finish := bid(item)
					fmt.Println("received bid", finish)
					bidsCh <- finish
				} else {
					awaitServerWorkCh <- true
					// mata goroutine
					fmt.Println("killing goroutine")
					return
				}
			}
		}()
	}

	// aguarda o processamento de todos os servidores
	for i := 1; i <= nServers; i++ {
		<-awaitServerWorkCh
	}
	close(bidsCh)
	return bidsCh
}

func main() {
	nServers := 3
	handle(nServers)
}
