package proposer

import (
	"sync"
	"log"
	"fmt"
)

type Proposal struct {
	id uint64
	message string
}

func NewProposal(id uint64, msg string) *Proposal {
	return &Proposal {id: id, message: msg}
}

func (proposal *Proposal) GetID() uint64 {
	return proposal.id
}

func (proposal *Proposal) GetMessage() string {
	return proposal.message
}

type Proposer struct {
	StopSig chan bool

	name string
	RecvProp chan Proposal
	proposals []Proposal
	lock sync.RWMutex
}

var g_proposers []*Proposer

func New(name string) *Proposer {
	ins := Proposer{name: name}
	ins.StopSig = make(chan bool)
	ins.RecvProp = make(chan Proposal)
	g_proposers = append(g_proposers, &ins)
	return &ins
}

func (proposer *Proposer) Start() {
	log.Println(fmt.Sprintf("%s开始接受提案", proposer.name))
	go func() {
		for {
			if proposal, ok := <- proposer.RecvProp; ok {
				if proposal.id == 0 {
					break
				}
				if proposal.id < proposer.maxId() {
					log.Println(fmt.Sprintf("%s拒绝了一份提案，编号为%d", proposer.name, proposal.id))
					continue
				}
				proposer.lock.Lock()
				proposer.proposals = append(proposer.proposals, proposal)
				log.Println(fmt.Sprintf("%s接受到一份提案，编号为%d", proposer.name, proposal.id))
				proposer.lock.Unlock()
			}
		}
		log.Println(fmt.Sprintf("%s停止接受提案", proposer.name))
		proposer.StopSig <- true
	}()
}

func (proposer *Proposer) Stop() {
	proposer.RecvProp <- Proposal{id: 0}
	<- proposer.StopSig
}

func (proposer *Proposer) maxId() uint64 {
	var ret uint64 = 0
	for _, prop := range proposer.proposals {
		if prop.id > ret {
			ret = prop.id
		}
	}
	return ret
}