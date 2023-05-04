package handler

import (
	"context"
	"encoding/json"
	"joerx/minecraft-cli/internal/service/rcon"
	"net/http"
)

type CommandFunc func(context.Context, rcon.CommandInput) (rcon.CommandOutput, error)

type CommandHandler struct {
	cmd CommandFunc
}

func NewCommandHandler(s *rcon.Service) *CommandHandler {
	return &CommandHandler{cmd: s.Command}
}

func (ch *CommandHandler) RunCommand(w http.ResponseWriter, r *http.Request) error {
	var input rcon.CommandInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return err
	}

	output, err := ch.cmd(r.Context(), input)
	if err != nil {
		return err
	}

	serveJSON(w, output, http.StatusAccepted)
	return nil
}
