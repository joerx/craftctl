package handler

import (
	"embed"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

// go:embed world
var worldFS embed.FS

type testRCon struct {
	t *testing.T
}

func (m testRCon) Command(cmd string) error {
	m.t.Logf("Received command '%s'", cmd)
	return nil
}

func TestHandleBackup(t *testing.T) {
	cm := &testRCon{t}
	handler := NewBackup(cm, worldFS)

	req := httptest.NewRequest(http.MethodPost, "/", nil)
	w := httptest.NewRecorder()

	handler.ServeHTTP(w, req)

	var d map[string]string
	resp := w.Result()

	if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
		t.Fatal(err)
	}

	wantCode := http.StatusOK
	if wantCode != resp.StatusCode {
		t.Errorf("want status %d, got %d", wantCode, resp.StatusCode)
	}

	wantSum := "76cdb2bad9582d23c1f6f4d868218d6c"
	gotSum := d["md5"]
	if wantSum != gotSum {
		t.Errorf("want md5 sum %s, got %s", wantSum, gotSum)
	}
}
