package main

import (
	"fmt"
	"time"
	"math/rand"
	"proposal"
	"log"
	"member"
)

func testAsync(propsr *member.Proposer, stopSig chan bool) {
	for i := 1; i < 5; i ++ {
		propslId := rand.Intn(100)
		if propslId == 0 {
			propslId = 1
		}
		proposal.GetPropslPool().SendProposal(
			*proposal.New(int64(propslId), fmt.Sprintf("消息%d", i)),
		)
		rand.Seed(time.Now().UTC().UnixNano())
		time.Sleep(time.Duration(rand.Intn(10)) * time.Second)
	}
	stopSig <- true
}

func main() {
	proposer1 := member.NewProposer("proposer1")
	acceptor1 := member.NewAcceptor("acceptor1")
	acceptor2 := member.NewAcceptor("acceptor2")
	acceptor3 := member.NewAcceptor("acceptor3")
	stopSig := make(chan bool, 3)
	go testAsync(proposer1, stopSig)
	go testAsync(proposer1, stopSig)
	go testAsync(proposer1, stopSig)
	<- stopSig
	<- stopSig
	<- stopSig
	proposal.GetPropslPool().Stop()
	log.Println(acceptor1.GetName(), acceptor1.ListProposals())
	log.Println(acceptor2.GetName(), acceptor2.ListProposals())
	log.Println(acceptor3.GetName(), acceptor3.ListProposals())
}
