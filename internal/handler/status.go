package handler

import (
	"joerx/minecraft-cli/internal/handler/task"
	"net/http"
)

type statusHandler struct {
	mon task.Monitor
}

func NewStatus(m task.Monitor) http.Handler {
	return &statusHandler{m}
}

func (h *statusHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	state, err := h.mon.GetState(req.Context())
	if err != nil {
		serveJSONError(w, err, http.StatusInternalServerError)
		return
	}
	serveJSON(w, state, http.StatusOK)
}
