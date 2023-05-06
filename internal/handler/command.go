package handler

import (
	"context"
	"encoding/json"
	"joerx/minecraft-cli/internal/api/rcon"
	"net/http"
)

type commandFunc func(context.Context, rcon.CommandInput) (rcon.CommandOutput, error)

type commandHandler struct {
	cmd commandFunc
}

// Alternative way to create a handler without using a closure - we still use our custom
// handler function for error handling and convert our custom ServeHTTP into the proper http.Handler type
func Command(fn commandFunc) http.Handler {
	ch := &commandHandler{cmd: fn}
	return h(ch.ServeHTTP)
}

func (ch *commandHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) error {
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
