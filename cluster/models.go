package cluster

//go:generate easyjson --all

type RegisterRequest struct {
	Service string
	ID      string
}

func (n RegisterRequest) NodeID() string    { return n.ID }
func (n RegisterRequest) ServiceID() string { return n.Service }
