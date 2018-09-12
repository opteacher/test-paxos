package member

import (
	"log"
	"fmt"
	"proposal"
)

type Proposer struct {
	name string
	sentProps []proposal.Proposal
}

var g_proposers []*Proposer

func NewProposer(name string) *Proposer {
	ins := &Proposer{name: name}
	ins.sentProps = []proposal.Proposal {}
	g_proposers = append(g_proposers, ins)
	ins.start()
	return ins
}

func (proposer *Proposer) EmitAccept(id int64) {
	var propsl proposal.Proposal
	for _, prop := range proposer.sentProps {
		if prop.GetID() == id {
			propsl = prop
		}
	}
	if propsl.GetID() == 0 {
		return
	}
	propsl.SetStatus("accept")
	acceptors := ShuffleListAcceptors()
	for _, acc := range acceptors {
		recv := acc.Send(propsl)
		log.Println(fmt.Sprintf("proposer:%s将编号为%d的提案发送给acceptor:%s，收到：%s",
			proposer.name, id, acc.GetName(), recv))
	}
}

func (proposer *Proposer) start() {
	log.Println(fmt.Sprintf("%s开始接受提案", proposer.name))
	go func() {
		for {
			propPool := proposal.GetPropslPool()
			propPool.Lock.RLock()
			propsl := propPool.TakeProposal()
			propPool.Lock.RUnlock()

			isHandled := false
			for _, prop := range proposer.sentProps {
				if prop.GetID() == propsl.GetID() {
					isHandled = true
					break
				}
			}
			if isHandled {
				continue
			}
			if propsl.GetID() == 0 {
				continue
			}
			if propsl.GetID() == -1 {
				break
			}

			log.Println(fmt.Sprintf("%s获取了一份编号为%d的提案", proposer.name, propsl.GetID()))
			propsl.SetReceiver(proposer.name)
			proposer.sentProps = append(proposer.sentProps, propsl)
			acceptors := ShuffleListAcceptors()
			for _, acc := range acceptors {
				recv := acc.Send(propsl)
				log.Println(fmt.Sprintf("编号为%d的提案收到来自%s的信息：%s",
					propsl.GetID(), acc.GetName(), recv))
			}
		}
	}()
}

func findProposer(name string) *Proposer {
	for _, proposer := range g_proposers {
		if proposer.name == name {
			return proposer
		}
	}
	return nil
}