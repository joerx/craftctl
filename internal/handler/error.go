package handler

import (
	"log"
	"net/http"
)

type errorResponse struct {
	Error string `json:"error"`
}

func serveJSONError(w http.ResponseWriter, err error, statusCode int) {
	log.Println(err)
	serveJSON(w, errorResponse{err.Error()}, statusCode)
}
