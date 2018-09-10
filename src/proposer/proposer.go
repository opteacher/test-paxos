package proposer

import (
	"log"
	"fmt"
	"acceptor"
	"proposal"
)

type Proposer struct {
	StopSig chan bool
	name string
	RecvProp chan proposal.Proposal
}

func New(name string) *Proposer {
	ins := Proposer{name: name}
	ins.StopSig = make(chan bool)
	ins.RecvProp = make(chan proposal.Proposal)
	return &ins
}

func (proposer *Proposer) Start() {
	log.Println(fmt.Sprintf("%s开始接受提案", proposer.name))
	go func() {
		for {
			if propsl, ok := <- proposer.RecvProp; ok {
				if propsl.GetID() == 0 {
					break
				}
				log.Println(fmt.Sprintf("%s接受到一份提案，编号为%d", proposer.name, propsl.GetID()))
				acceptors := acceptor.ShuffleListAcceptors()
				for _, acc := range acceptors {
					recv := acc.Send(propsl)
					log.Println(fmt.Sprintf("编号为%d的提案收到来自%s的信息：%s",
						propsl.GetID(), acc.GetName(), recv))
				}
			}
		}
		log.Println(fmt.Sprintf("%s停止接受提案", proposer.name))
		proposer.StopSig <- true
	}()
}

func (proposer *Proposer) Stop() {
	proposer.RecvProp <- *proposal.New(0, "")
	<- proposer.StopSig
}