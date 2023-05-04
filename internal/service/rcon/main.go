package rcon

import "context"

type RCon interface {
	Command(ctx context.Context, cmd string) (string, error)
}

type Service struct {
	rcon RCon
}

func NewService(rc RCon) *Service {
	return &Service{rcon: rc}
}

type CommandInput struct {
	Command string `json:"command"`
}

type CommandOutput struct {
	Response string `json:"response"`
}

func (s *Service) Command(ctx context.Context, input CommandInput) (CommandOutput, error) {
	resp, err := s.rcon.Command(ctx, input.Command)
	if err != nil {
		return CommandOutput{}, err
	}
	return CommandOutput{Response: resp}, nil
}
