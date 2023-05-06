package handler

import "net/http"

// Custom http.Handler that takes care of error handling and implements the http.Handler interface
type h func(w http.ResponseWriter, req *http.Request) error

func (fn h) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	// From https://chidiwilliams.com/post/writing-cleaner-go-web-servers/
	err := fn(w, req)
	if e, ok := err.(httpError); ok {
		serveJSONError(w, err, e.code)
		return
	}

	if err != nil {
		serveJSONError(w, err, http.StatusInternalServerError)
		return
	}
}

type httpError struct {
	code int
	err  error
}

func badRequest(err error) httpError {
	return httpError{code: http.StatusBadRequest, err: err}
}

func (ie httpError) Error() string {
	return ie.err.Error()
}
