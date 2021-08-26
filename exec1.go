package main

import (
	"math/rand"
	"time"
)

func gateway(nreplicas int) int32 {
	respCh := make(chan int32)
	for i := 0; i < nreplicas; i++ {
		go request(respCh)
	}
	first := <-respCh
	return first
}

func request(out chan int32) {
	rand.Seed(time.Now().UnixNano())
	number := rand.Int31n(30) + 1
	println("sleeping seconds", number)
	time.Sleep(time.Duration(number) * time.Second)
	out <- number
}

func main() {
	value := gateway(3)
	println("gateway value", value)
}
