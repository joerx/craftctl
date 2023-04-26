package handler

import (
	"joerx/minecraft-cli/internal/handler/task"
	"net/http"
)

type stopHandler struct {
	stopper task.Stopper
}

func NewStop(sr task.Stopper) http.Handler {
	return &stopHandler{sr}
}

func (h *stopHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	state, err := h.stopper.Stop(req.Context())
	if err != nil {
		serveJSONError(w, err, http.StatusInternalServerError)
		return
	}
	serveJSON(w, statusResponse{state}, http.StatusOK)
}
