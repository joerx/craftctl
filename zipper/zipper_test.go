package zipper

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"testing"
	"testing/fstest"
)

func TestZipFS(t *testing.T) {
	want := "0ce9b4eb5cf5570a20087bf177f4d8d4"
	fsys := fstest.MapFS{
		"file.go":                {Data: []byte("fmt.Println(\"Hello gopher\")")},
		"subfolder/subfolder.go": {Data: []byte("fmt.Println(\"package subfolder\")")},
		"subfolder2/another.go":  {Data: []byte("fmt.Println(\"package subfolder2\")")},
		"subfolder2/file.go":     {Data: []byte("fmt.Println(\"package subfolder2\")")},
	}

	w := new(bytes.Buffer)
	if err := ZipFS(fsys, w); err != nil {
		t.Fatal(err)
	}

	got := fmt.Sprintf("%x", md5.Sum(w.Bytes()))
	if got != want {
		t.Errorf("want %s, got %s", want, got)
	}
}
