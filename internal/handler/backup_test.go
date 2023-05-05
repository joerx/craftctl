package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"joerx/minecraft-cli/internal/service/backup"
	"net/http"
	"net/http/httptest"
	"testing"
)

// go:embed testdata/world
// var worldFS embed.FS
// var worldSum string = "76cdb2bad9582d23c1f6f4d868218d6c"

type testBackupService struct {
}

func (tb *testBackupService) Create(ctx context.Context, in backup.CreateBackupInput) (backup.CreateBackupOutput, error) {
	return backup.CreateBackupOutput{}, nil
}

func (tb *testBackupService) List(ctx context.Context) (backup.ListBackupOutput, error) {
	return backup.ListBackupOutput{}, nil
}

func (tb *testBackupService) Restore(ctx context.Context, in backup.RestoreBackupInput) (backup.RestoreBackupOutput, error) {
	return backup.RestoreBackupOutput{}, nil
}

func TestCreateBackup(t *testing.T) {
	handler := &BackupHandler{svc: &testBackupService{}}

	in := backup.CreateBackupInput{Key: "foo"}
	bdy := bytes.NewBuffer([]byte{})
	if err := json.NewEncoder(bdy).Encode(in); err != nil {
		t.Fatal(err)
	}

	req := httptest.NewRequest(http.MethodPost, "/", bdy)
	w := httptest.NewRecorder()

	if err := handler.Create(w, req); err != nil {
		t.Fatal(err)
	}

	var o backup.CreateBackupOutput
	resp := w.Result()

	if err := json.NewDecoder(resp.Body).Decode(&o); err != nil {
		t.Fatal(err)
	}

	wantCode := http.StatusOK
	if wantCode != resp.StatusCode {
		t.Errorf("want status %d, got %d", wantCode, resp.StatusCode)
	}
}
