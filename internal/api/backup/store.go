package backup

import (
	"context"
	"io"
)

type ObjectInfo struct {
	Location string `json:"location"`
	Key      string `json:"key"`
}

type Putter interface {
	Put(ctx context.Context, key string, r io.Reader) (ObjectInfo, error)
}

type Getter interface {
	Get(ctx context.Context, key string, w io.Writer) (ObjectInfo, error)
}

type Lister interface {
	List(ctx context.Context) ([]ObjectInfo, error)
}

type Store interface {
	Putter
	Lister
	// Getter
}
