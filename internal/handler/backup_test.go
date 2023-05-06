package handler

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"joerx/minecraft-cli/internal/api/backup"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"
)

func TestCreateBackup(t *testing.T) {
	testCases := []struct {
		in           backup.CreateBackupInput
		out          backup.CreateBackupOutput
		err          error
		wantCode     int
		wantResponse map[string]string // Response may be either an error response or a CreateBackupOutput
	}{
		{
			in:           backup.CreateBackupInput{Key: "foo"},
			out:          backup.CreateBackupOutput{MD5: "abc123", ObjectInfo: backup.ObjectInfo{Location: "test://foo", Key: "foo"}},
			err:          nil,
			wantCode:     http.StatusOK,
			wantResponse: map[string]string{"md5": "abc123", "location": "test://foo", "key": "foo"},
		},
		{
			in:           backup.CreateBackupInput{},
			out:          backup.CreateBackupOutput{},
			err:          backup.InputError("invalid input"),
			wantCode:     http.StatusBadRequest,
			wantResponse: map[string]string{"error": "invalid input"},
		},
		{
			in:           backup.CreateBackupInput{},
			out:          backup.CreateBackupOutput{},
			err:          fmt.Errorf("the server made a boo boo"),
			wantCode:     http.StatusInternalServerError,
			wantResponse: map[string]string{"error": "the server made a boo boo"},
		},
	}

	for _, tc := range testCases {
		// Each handler only requires the service function with the matching interface instead of a "fat" service interface
		// This makes handlers easy to test since we only mock the matching service function used by the handler
		// Using closures we can even do this inside a tabular test loop - no need for complex mocking
		h := CreateBackup(func(ctx context.Context, in backup.CreateBackupInput) (backup.CreateBackupOutput, error) {
			return tc.out, tc.err
		})

		bdy := bytes.NewBuffer([]byte{})
		if err := json.NewEncoder(bdy).Encode(tc.in); err != nil {
			t.Fatal(err)
		}

		req := httptest.NewRequest(http.MethodPost, "/", bdy)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		resp := w.Result()

		if tc.wantCode != resp.StatusCode {
			t.Errorf("want status %d, got %d", tc.wantCode, resp.StatusCode)
		}

		var gotResponse map[string]string
		if err := json.NewDecoder(resp.Body).Decode(&gotResponse); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(tc.wantResponse, gotResponse) {
			t.Errorf("want response %#v but got %#v", tc.wantResponse, gotResponse)
		}
	}
}

func TestListBackups(t *testing.T) {
	testCases := []struct {
		out          backup.ListBackupOutput
		err          error
		wantCode     int
		wantResponse map[string]interface{}
	}{
		{
			out:      backup.ListBackupOutput{Backups: []backup.ObjectInfo{{Location: "test://foo", Key: "foo"}}},
			err:      nil,
			wantCode: http.StatusOK,
			wantResponse: map[string]interface{}{
				"backups": []interface{}{
					map[string]interface{}{
						"location": "test://foo",
						"key":      "foo",
					},
				},
			},
		},
		{
			out:      backup.ListBackupOutput{},
			err:      fmt.Errorf("whoopsie"),
			wantCode: http.StatusInternalServerError,
			wantResponse: map[string]interface{}{
				"error": "whoopsie",
			},
		},
	}

	for _, tc := range testCases {
		h := ListBackups(func(ctx context.Context) (backup.ListBackupOutput, error) {
			return tc.out, tc.err
		})

		req := httptest.NewRequest(http.MethodGet, "/", nil)
		w := httptest.NewRecorder()

		h.ServeHTTP(w, req)
		resp := w.Result()

		if tc.wantCode != resp.StatusCode {
			t.Errorf("want status %d, got %d", tc.wantCode, resp.StatusCode)
		}

		var gotResponse map[string]interface{}
		if err := json.NewDecoder(resp.Body).Decode(&gotResponse); err != nil {
			t.Fatal(err)
		}
		if !reflect.DeepEqual(tc.wantResponse, gotResponse) {
			t.Errorf("want response %#v but got %#v", tc.wantResponse, gotResponse)
		}
	}
}
