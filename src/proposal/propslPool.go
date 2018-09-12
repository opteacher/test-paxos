package proposal

import (
	"sync"
	"log"
	"fmt"
)

type PropslPool struct {
	proposals []Proposal
	recvProp chan Proposal
	stopSig chan bool
	Lock sync.RWMutex
}

var g_propslPool *PropslPool

func GetPropslPool() *PropslPool {
	if g_propslPool == nil {
		g_propslPool = &PropslPool {}
		g_propslPool.recvProp = make(chan Proposal)
		g_propslPool.stopSig = make(chan bool)
		g_propslPool.start()
	}
	return g_propslPool
}

func (propslPool *PropslPool) start() {
	go func() {
		for {
			if propsl, ok := <- propslPool.recvProp; ok {
				if propsl.GetID() == -1 {
					break
				}
				propslPool.Lock.Lock()
				propslPool.proposals = append(propslPool.proposals, propsl)
				propslPool.Lock.Unlock()
				log.Println(fmt.Sprintf("收到一份提案，编号为%d", propsl.GetID()))
			}
		}
		log.Println("提案池停止接受提案")
		propslPool.stopSig <- true
	}()
}

func (propslPool *PropslPool) Stop() {
	propslPool.recvProp <- Proposal{id: -1}
	<- propslPool.stopSig
}

func (propslPool *PropslPool) TakeProposal() Proposal {
	propslNum := len(propslPool.proposals)
	if propslNum != 0 {
		ret := propslPool.proposals[propslNum - 1]
		return ret
	} else {
		return Proposal {}
	}
}

func (propslPool *PropslPool) SendProposal(proposal Proposal) {
	propslPool.recvProp <- proposal
}