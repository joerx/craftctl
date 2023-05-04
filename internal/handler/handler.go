package handler

import "net/http"

type Handler func(w http.ResponseWriter, req *http.Request) error

func (h Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// From https://chidiwilliams.com/post/writing-cleaner-go-web-servers/
	if err := h(w, req); err != nil {
		serveJSONError(w, err, http.StatusInternalServerError)
	}
}
