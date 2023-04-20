package zipper

import (
	"archive/zip"
	"fmt"
	"io"
	"io/fs"
)

func ZipFS(fsys fs.FS, w io.Writer) error {
	zw := zip.NewWriter(w)

	err := fs.WalkDir(fsys, ".", func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return fmt.Errorf("error reading entry %s: %v", path, err)
		}

		if d.IsDir() {
			return nil
		}

		zf, err := zw.Create(path)
		if err != nil {
			return fmt.Errorf("error creating entry %s: %v", path, err)
		}

		fd, err := fsys.Open(path)
		if err != nil {
			return fmt.Errorf("error opening %s: %v", path, err)
		}

		if _, err := io.Copy(zf, fd); err != nil {
			return fmt.Errorf("error copying %s to output: %v", path, err)
		}

		return nil
	})

	if err != nil {
		return err
	}

	if err := zw.Close(); err != nil {
		return err
	}

	return nil
}
