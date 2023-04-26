package handler

import (
	"context"
	"encoding/json"
	"fmt"
	"joerx/minecraft-cli/internal/handler/task"
	"net/http"
	"net/http/httptest"
	"testing"
)

type testMonitor struct {
	result   task.State
	err      error
	wantCode int
}

func (m *testMonitor) GetState(ctx context.Context) (task.State, error) {
	return m.result, m.err
}

func TestHandleStatus(t *testing.T) {
	testCases := []*testMonitor{
		{task.State{State: task.StateActive}, nil, http.StatusOK},
		{task.State{}, fmt.Errorf("Something went wrong"), http.StatusInternalServerError},
	}

	for _, tc := range testCases {
		handler := NewStatus(tc)
		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		handler.ServeHTTP(w, req)

		resp := w.Result()

		if resp.StatusCode != tc.wantCode {
			t.Errorf("want status %d, got %d", tc.wantCode, resp.StatusCode)
		}

		wantType := "application/json"
		gotType := resp.Header.Get("Content-type")
		if gotType != wantType {
			t.Errorf("want content-type '%s', got '%s'", wantType, gotType)
		}

		var d map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
			t.Fatal(err)
		}

		wantState := tc.result.State
		gotState := d["state"]
		if gotState != wantState {
			t.Errorf("want task state '%s', got '%s'", wantState, gotState)
		}
	}
}
