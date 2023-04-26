package frontend

import (
	"embed"
	"io"
	"io/fs"
	"log"
	"mime"
	"net/http"
	"os"
	"path"
	"strings"
)

//go:embed assets
//go:embed index.html
var staticFS embed.FS

type fsHandler struct {
	fs fs.FS
}

func handleIndex(w http.ResponseWriter, r *http.Request) {
	// See https://pkg.go.dev/net/http#ServeMux
	if r.URL.Path != "/" && r.URL.Path != "/index.html" {
		http.NotFound(w, r)
		return
	}

	data, _ := staticFS.Open("index.html")
	io.Copy(w, data)
}

func (h fsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Printf("%s %s", r.Method, r.URL.Path)

	pt := strings.TrimPrefix(r.URL.Path, "/")
	f, err := h.fs.Open(pt)
	switch {
	case os.IsNotExist(err):
		http.Error(w, "File not found", http.StatusNotFound)
		return
	case err != nil:
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	ext := path.Ext(pt)
	ct := mime.TypeByExtension(ext)

	w.Header().Add("Content-Type", ct)
	w.WriteHeader(http.StatusOK)

	io.Copy(w, f)
}

func New() http.Handler {
	// return
	mux := http.NewServeMux()
	mux.HandleFunc("/", handleIndex)
	mux.Handle("/assets/", &fsHandler{staticFS})
	return mux
}
