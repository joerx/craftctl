package handler

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

type statusResponse struct {
	Status string `json:"status"`
}

type errorResponse struct {
	Error string `json:"error"`
}

func serveJSON(w http.ResponseWriter, payload interface{}, statusCode int) {
	buf := bytes.NewBuffer([]byte{})

	if err := json.NewEncoder(buf).Encode(payload); err != nil {
		http.Error(w, fmt.Sprintf("failed to serialize response - %v", err), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-type", "application/json")
	w.WriteHeader(statusCode)
	io.Copy(w, buf)
}

func serveJSONError(w http.ResponseWriter, err error, statusCode int) {
	serveJSON(w, errorResponse{err.Error()}, statusCode)
}
