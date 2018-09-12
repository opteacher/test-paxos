package member

import (
	"log"
	"fmt"
	"proposal"
	"math/rand"
)

type Acceptor struct {
	name string
	proposals []proposal.Proposal
}

var g_acceptors []*Acceptor

func NewAcceptor(name string) *Acceptor {
	ins := Acceptor{name:name}
	g_acceptors = append(g_acceptors, &ins)
	return &ins
}

func ShuffleListAcceptors() []*Acceptor {
	ret := g_acceptors
	rand.Shuffle(len(ret), func(i, j int) {
		ret[i], ret[j] = ret[j], ret[i]
	})
	return ret
}

func (acceptor *Acceptor) Send(proposal proposal.Proposal) string {
	log.Println(fmt.Sprintf("%s接受到一份提案，编号为%d", acceptor.name, proposal.GetID()))
	if proposal.GetStatus() == "accept" {
		log.Println(fmt.Sprintf("编号为%d提案进入accept状态", proposal.GetID()))
		var maxId int64 = 0
		index := 0
		for i, prop := range acceptor.proposals {
			if prop.GetID() == proposal.GetID() {
				acceptor.proposals[i].SetStatus("accept")
				index = i
			}
			if prop.GetStatus() == "accept" || prop.GetStatus() == "passed" {
				if prop.GetID() > maxId {
					maxId = prop.GetID()
				}
			}
		}
		if proposal.GetID() > maxId && index < len(acceptor.proposals) {
			log.Println(fmt.Sprintf("acceptor:%s未接收过编号大于%d的提案，所以提案通过",
				acceptor.name, proposal.GetID()))
			acceptor.proposals[index].SetStatus("passed")
			return fmt.Sprintf("编号为%d的提案通过", proposal.GetID())
		}
		return proposal.GetMessage()
	}
	bestProposal := acceptor.bestProposal()
	if bestProposal != nil && bestProposal.GetID() > proposal.GetID() {
		return fmt.Sprintf("%s的最大编号为%d，所以拒绝了编号为%d的提案",
			acceptor.name, bestProposal.GetID(), proposal.GetID())
	}
	propslNum := 0
	for _, acc := range g_acceptors {
		if acc.hasProposal(proposal.GetID()) {
			propslNum++
		}
	}
	if propslNum > (len(g_acceptors) >> 1) {
		log.Println(fmt.Sprintf("编号为%d的提案经手%s的时候，进入commit状态", proposal.GetID(), acceptor.name))
		proposal.SetStatus("commit")
		findProposer(proposal.GetReceiver()).EmitAccept(proposal.GetID())
	}
	acceptor.proposals = append(acceptor.proposals, proposal)
	if bestProposal != nil {
		return bestProposal.GetMessage()
	} else {
		return ""
	}
}

func (acceptor *Acceptor) bestProposal() *proposal.Proposal {
	if len(acceptor.proposals) == 0 {
		return nil
	} else {
		bestProposal := &acceptor.proposals[0]
		for _, propsl := range acceptor.proposals {
			if propsl.GetID() > bestProposal.GetID() {
				bestProposal = &propsl
			}
		}
		return bestProposal
	}
}

func (acceptor *Acceptor) hasProposal(id int64) bool {
	for _, propsl := range acceptor.proposals {
		if propsl.GetID() == id {
			return true
		}
	}
	return false
}

func (acceptor *Acceptor) GetName() string {
	return acceptor.name
}

func (acceptor *Acceptor) ListProposals() []proposal.Proposal {
	return acceptor.proposals
}