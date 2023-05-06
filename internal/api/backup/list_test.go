package backup

import (
	"context"
	"fmt"
	"io"
	"testing"
)

type testStore struct {
	backups []ObjectInfo
	err     error
}

func (ts *testStore) Put(ctx context.Context, key string, r io.Reader) (ObjectInfo, error) {
	return ObjectInfo{}, ts.err
}

func (ts *testStore) List(ctx context.Context) ([]ObjectInfo, error) {
	return ts.backups, ts.err
}

func (ts *testStore) Get(ctx context.Context, key string, w io.WriterAt) error {
	return nil
}

func TestList(t *testing.T) {
	testCases := []struct {
		wantErr    error
		wantResult ListBackupOutput
	}{
		{
			wantErr: nil,
			wantResult: ListBackupOutput{
				Backups: []ObjectInfo{
					{
						Location: "test://foo",
						Key:      "bla",
					},
				},
			},
		},
		{
			wantErr:    fmt.Errorf("kaboom!"),
			wantResult: ListBackupOutput{},
		},
	}

	for _, tc := range testCases {
		s := &backupService{store: &testStore{
			backups: tc.wantResult.Backups,
			err:     tc.wantErr,
		}}

		o, err := s.List(context.Background())

		if err != tc.wantErr {
			t.Errorf("want error %#v but got %#v", tc.wantErr, err)
		}

		if len(tc.wantResult.Backups) != len(o.Backups) {
			t.Errorf("wanted %d backup(s), but got %d", len(tc.wantResult.Backups), len(o.Backups))
		}
	}
}
