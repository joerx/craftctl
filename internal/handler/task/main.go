package task

import "context"

const (
	StateActive   = "active"
	StateInactive = "inactive"
	StateUnknown  = "unknown"
)

type State struct {
	State       string `json:"state"`
	StateDetail string `json:"state-detail"`
	Name        string `json:"name"`
}

type Monitor interface {
	GetState(ctx context.Context) (State, error)
}

type Stopper interface {
	Stop(ctx context.Context) (string, error)
}

type Starter interface {
	Start(ctx context.Context) (string, error)
}
