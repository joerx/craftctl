package handler

import (
	"context"
	"encoding/json"
	"net/http"
)

const StateActive = "active"
const StateInactive = "inactive"
const StateUnknown = "unknown"

type Starter interface {
	Start(ctx context.Context) error
}

type Stopper interface {
	Stop(ctx context.Context) error
}

type State struct {
	State       string `json:"state"`
	StateDetail string `json:"state-detail"`
	Name        string `json:"name"`
}

type Monitor interface {
	GetState(ctx context.Context) (State, error)
}

type startHandler struct {
	starter Starter
}

type stopHandler struct {
	stopper Stopper
}

type statusHandler struct {
	mon Monitor
}

func NewStart(sr Starter) http.Handler {
	return &startHandler{sr}
}

func NewStop(sr Stopper) http.Handler {
	return &stopHandler{sr}
}

func NewStatus(m Monitor) http.Handler {
	return &statusHandler{m}
}

func (h *startHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "TODO", http.StatusNotImplemented)
}

func (h *stopHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	http.Error(w, "TODO", http.StatusNotImplemented)
}

func (h *statusHandler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	state, err := h.mon.GetState(req.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")

	json.NewEncoder(w).Encode(state)
}
