package main

import (
	"proposer"
	"fmt"
	"time"
	"math/rand"
	"acceptor"
	"proposal"
	"log"
)

func testAsync(propsr *proposer.Proposer, stopSig chan bool) {
	for i := 1; i < 10; i ++ {
		propslId := rand.Intn(100)
		if propslId == 0 {
			propslId = 1
		}
		propsr.RecvProp <- *proposal.New(uint64(propslId), fmt.Sprintf("消息%d", i))
		rand.Seed(time.Now().UTC().UnixNano())
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
	stopSig <- true
}

func main() {
	proposer1 := proposer.New("proposer1")
	acceptor1 := acceptor.New("acceptor1")
	acceptor2 := acceptor.New("acceptor2")
	acceptor3 := acceptor.New("acceptor3")
	proposer1.Start()
	stopSig := make(chan bool, 3)
	go testAsync(proposer1, stopSig)
	go testAsync(proposer1, stopSig)
	go testAsync(proposer1, stopSig)
	<- stopSig
	<- stopSig
	<- stopSig
	proposer1.Stop()
	log.Println(acceptor1.GetName(), acceptor1.ListProposals())
	log.Println(acceptor2.GetName(), acceptor2.ListProposals())
	log.Println(acceptor3.GetName(), acceptor3.ListProposals())
}
