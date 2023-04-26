package handler

import (
	"joerx/minecraft-cli/internal/handler/task"
	"net/http"
)

type startHandler struct {
	starter task.Starter
}

func NewStart(sr task.Starter) http.Handler {
	return &startHandler{sr}
}

func (h *startHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	state, err := h.starter.Start(req.Context())
	if err != nil {
		serveJSONError(w, err, http.StatusInternalServerError)
		return
	}
	serveJSON(w, statusResponse{state}, http.StatusOK)
}
