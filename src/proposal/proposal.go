package proposal

type Proposal struct {
	id uint64
	message string
	status string
}

func New(id uint64, msg string) *Proposal {
	return &Proposal {id: id, message: msg, status: "prepare"}
}

func (proposal *Proposal) GetID() uint64 {
	return proposal.id
}

func (proposal *Proposal) GetMessage() string {
	return proposal.message
}

func (proposal *Proposal) SetStatus(status string) {
	proposal.status = status
}