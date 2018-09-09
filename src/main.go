package main

import (
	"proposer"
	"fmt"
	"time"
	"math/rand"
)

func testAsync(propsr *proposer.Proposer, stopSig chan bool) {
	for i := 1; i < 10; i ++ {
		propsr.RecvProp <- *proposer.NewProposal(uint64(i), fmt.Sprintf("消息%d", i))
		rand.Seed(time.Now().UTC().UnixNano())
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
	stopSig <- true
}

func main() {
	proposer1 := proposer.New("proposer1")
	proposer1.Start()
	stopSig := make(chan bool, 3)
	go testAsync(proposer1, stopSig)
	go testAsync(proposer1, stopSig)
	go testAsync(proposer1, stopSig)
	<- stopSig
	<- stopSig
	<- stopSig
	proposer1.Stop()
}
