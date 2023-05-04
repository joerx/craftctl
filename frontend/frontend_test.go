package frontend

import (
	"crypto/md5"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
	"testing/fstest"
)

func TestFSHandler(t *testing.T) {
	fsys := fstest.MapFS{
		"index.html":     {Data: []byte("<h1>Hello World<h1>")},
		"css/styles.css": {Data: []byte("body {background-image: url(\"unicorn.png\")}")},
		"js/script.js":   {Data: []byte("console.log('Hello World')")},
	}

	tests := []struct {
		path     string
		wantCode int
		wantMime string
		wantSum  [16]byte
	}{
		{"index.html", http.StatusOK, "text/html", md5.Sum(fsys["index.html"].Data)},
		{"foo.html", http.StatusNotFound, "text/plain", [16]byte{}},
	}

	handler := &fsHandler{fs: fsys}

	for _, tc := range tests {
		req := httptest.NewRequest("GET", fmt.Sprintf("http://example.com/%s", tc.path), nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			t.Fatal(err)
		}

		if resp.StatusCode != tc.wantCode {
			t.Errorf("expected code %d but got %d", tc.wantCode, resp.StatusCode)
		}

		gotMime := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(gotMime, tc.wantMime) {
			t.Errorf("want content type to start with '%s' but got '%s'", tc.wantMime, gotMime)
		}

		gotSum := md5.Sum(body)

		// Only check content if wanted checksum is not empty
		if tc.wantSum != [16]byte{} && gotSum != tc.wantSum {
			t.Errorf("want content md5 sum %x but got %x", tc.wantSum, gotSum)
		}
	}
}

func TestFrontendHandler(t *testing.T) {
	handler := New()

	tests := []struct {
		path     string
		wantMime string
		wantCode int
	}{
		{"/", "text/html", http.StatusOK},
		{"/index.html", "text/html", http.StatusOK},
		{"/foo.html", "text/plain", http.StatusNotFound},
		{"/assets/js/script.js", "text/javascript", http.StatusOK},
		{"/assets/css/style.css", "text/css", http.StatusOK},
	}

	for _, tc := range tests {
		req := httptest.NewRequest(http.MethodGet, fmt.Sprintf("http://example.com%s", tc.path), nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()
		if resp.StatusCode != tc.wantCode {
			t.Errorf("want code %d for path %s but got %d", tc.wantCode, tc.path, resp.StatusCode)
		}

		gotMime := resp.Header.Get("Content-Type")
		if !strings.HasPrefix(gotMime, tc.wantMime) {
			t.Errorf("want mime starting with '%s' for path %s but got '%s'", tc.wantMime, tc.path, gotMime)
		}
	}
}
