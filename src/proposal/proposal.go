package proposal

type Proposal struct {
	id int64
	message string
	status string
	receiver string
}

func New(id int64, msg string) *Proposal {
	return &Proposal {id: id, message: msg, status: "prepare"}
}

func (proposal *Proposal) GetID() int64 {
	return proposal.id
}

func (proposal *Proposal) GetMessage() string {
	return proposal.message
}

func (proposal *Proposal) SetStatus(status string) {
	proposal.status = status
}

func (proposal *Proposal) GetStatus() string {
	return proposal.status
}

func (proposal *Proposal) SetReceiver(name string) {
	proposal.receiver = name
}

func (proposal *Proposal) GetReceiver() string {
	return proposal.receiver
}