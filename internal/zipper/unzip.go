package zipper

import (
	"archive/zip"
	"io"
	"log"
	"os"
	"path"
)

func Unzip(path string, prefix string) error {
	zr, err := zip.OpenReader(path)
	if err != nil {
		return err
	}
	defer zr.Close()

	for _, f := range zr.File {
		if err := extractFile(f, prefix); err != nil {
			return err
		}
	}

	return nil
}

func extractFile(f *zip.File, prefix string) error {
	rc, err := f.Open()
	if err != nil {
		return err
	}
	defer rc.Close()

	dst := path.Join(prefix, f.Name)

	if err := os.MkdirAll(path.Dir(dst), 0755); err != nil {
		return err
	}

	of, err := os.Create(dst)
	if err != nil {
		return err
	}

	c, err := io.Copy(of, rc)
	if err != nil {
		return err
	}

	log.Printf("Extracted %s (%d bytes)", dst, c)
	return nil
}
