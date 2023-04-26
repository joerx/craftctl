package handler

import (
	"encoding/json"
	"log"
	"net/http"
)

type cmdRequest struct {
	Cmd string `json:"cmd"`
}

type Command struct {
	rcon RCon
}

func NewCommand(rc RCon) *Command {
	return &Command{rcon: rc}
}

func (ch *Command) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	var c cmdRequest
	if err := json.NewDecoder(r.Body).Decode(&c); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := ch.rcon.Command(c.Cmd); err != nil {
		log.Printf("Error sending command to server: %v", err)
		serveJSONError(w, err, http.StatusInternalServerError)
		return
	}

	serveJSON(w, statusResponse{Status: "accepted"}, http.StatusAccepted)
}
