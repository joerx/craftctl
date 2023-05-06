package instance

import "context"

const (
	StateActive   = "active"
	StateInactive = "inactive"
	StateUnknown  = "unknown"
)

type State struct {
	State  string
	Detail string
}

type Starter interface {
	Start(ctx context.Context) (State, error)
}

type Stopper interface {
	Stop(ctx context.Context) (State, error)
}

type Monitor interface {
	GetState(ctx context.Context) (State, error)
}

type Controller interface {
	Starter
	Stopper
	Monitor
}

type Service struct {
	controller Controller
}

type StartInstanceInput struct{}

type StopInstanceInput struct{}

type GetInstanceStateInput struct{}

type InstanceStateOutput struct {
	State State
}

func (s *Service) Start(ctx context.Context, input StartInstanceInput) (InstanceStateOutput, error) {
	return InstanceStateOutput{}, nil
}

func (s *Service) Stop(ctx context.Context, input StopInstanceInput) (InstanceStateOutput, error) {
	return InstanceStateOutput{}, nil
}

func (s *Service) GetState(ctx context.Context, input GetInstanceStateInput) (InstanceStateOutput, error) {
	return InstanceStateOutput{}, nil
}
