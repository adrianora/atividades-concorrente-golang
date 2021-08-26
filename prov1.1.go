package main

import (
	"math/rand"
	"time"
)

func receive(emptyCh chan int, receiveCh chan int, awaitCh chan int) {
	for {
		select {
		case tmp := <-receiveCh:
			println("out", tmp)
		default:
			<-awaitCh
		}
	}
}

func handle(limit int, emptyCh chan int, receiveCh chan int, awaitCh chan int) {
	for {
		time.Sleep(10 * time.Microsecond)
		rand.Seed(time.Now().UnixNano())
		number := rand.Intn(30)
		println("in ", number)
		receiveCh <- number
		if len(receiveCh) == limit {
			for i := 0; i < limit; i++ {
				awaitCh <- 1
			}
		}
	}
}

func main() {
	limit := 5
	emptyCh := make(chan int)
	receiveCh := make(chan int, limit)
	awaitCh := make(chan int)
	for i := 0; i < limit; i++ {
		go receive(emptyCh, receiveCh, awaitCh)
	}
	go handle(limit, emptyCh, receiveCh, awaitCh)
	for {
	}
}
