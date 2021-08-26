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
	sleepTime := rand.Intn(8)
	// na chegada de dois items simultaneamente em um canal
	// como o select escolhe qual caso tratar de forma aletória
	// o if garante que se o tempo de sleep for igual ao timeout
	// é incrementado em 2s para evitar a condição de corrida do select
	// está hardcoded para não alterar assinatura da função
	timeoutSecs := 5
	if sleepTime == timeoutSecs {
		sleepTime += 2
	}
	fmt.Printf("sleeping for %d seconds in item %+v\n", sleepTime, item)
	time.Sleep(time.Duration(sleepTime) * time.Second)

	return *bided
}

func handle(nServers int, timeoutSecs int) chan Bid {
	itemsCh := itemsStream()
	// tamanho do buffer é um número mágico que permite ser feito N tentativas
	bidsCh := make(chan Bid, 1000)
	// canal responsável por controlar a existência de todas as goroutines
	awaitServerWorkCh := make(chan bool, nServers)
	// quantos lances tiveram sucesso
	totalItems := len(itemsCh)
	countSuccessfulBids := 0

	for i := 1; i <= nServers; i++ {
		// cria uma goroutine para  cada servidor
		go func() {
			for {
				tickCh := time.Tick(time.Duration(timeoutSecs) * time.Second)
				item, notClosed := <-itemsCh
				if notClosed {
					timeoutCh := make(chan Bid)
					// cria uma nova goroutine que será responsável por
					// processar o lance. em seguida fico esperando o
					// processamento através do select. se o tick do
					// relógio bater antes do retorno da função bid, que
					// está sendo controlado através do canal timeoutCh,
					// é criado o bidFailed e adicionado no canal de saída
					go func() {
						finish := bid(item)
						timeoutCh <- finish
					}()
					select {
					case finish := <-timeoutCh:
						bidsCh <- finish
						fmt.Println("received bid", finish)
						// se todos os lances tiverem sucesso, fecha o canal
						countSuccessfulBids += 1
						if countSuccessfulBids == totalItems {
							close(itemsCh)
						}
					case <-tickCh:
						bidFailed := Bid{item, -1, true}
						bidsCh <- bidFailed
						fmt.Println("failed bid", item)
						// repõe o item no canal de items, tornando possível
						// ser feito novos lances
						itemsCh <- item
					}
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
	// no terminal sera mostrado que apenas uma goroutine será finalizada
	// para visualizar as demais goroutines serem finalizadas, retire o comentário
	// do Sleep na linha abaixo
	// time.Sleep(time.Duration(2) * time.Second)
	close(bidsCh)
	return bidsCh
}

func main() {
	nServers := 3
	timeoutSecs := 5
	handle(nServers, timeoutSecs)
}
