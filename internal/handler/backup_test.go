package handler

// // go:embed world
// var worldFS embed.FS
// var worldSum string = "76cdb2bad9582d23c1f6f4d868218d6c"

// type testRCon struct {
// 	t *testing.T
// }

// func (m testRCon) Command(cmd string) error {
// 	m.t.Logf("Received command '%s'", cmd)
// 	return nil
// }

// type testStore struct {
// 	prefix string
// }

// func (ts *testStore) Put(ctx context.Context, key string, r io.Reader) (ObjectInfo, error) {
// 	l := fmt.Sprintf("%s/%s", ts.prefix, key)
// 	return ObjectInfo{Location: l}, nil
// }

// func TestHandleBackup(t *testing.T) {
// 	prefix := "test://abc123"

// 	cm := &testRCon{t}
// 	ts := &testStore{prefix}
// 	handler := NewBackup(cm, worldFS, ts)

// 	testCases := []struct {
// 		wantCode     int
// 		wantSum      string
// 		wantLocation string
// 		key          string
// 	}{
// 		{http.StatusOK, worldSum, fmt.Sprintf("%s/bar.zip", prefix), "bar"},
// 		{http.StatusOK, worldSum, fmt.Sprintf("%s/foo-bar.zip", prefix), "foo bar"},
// 		{http.StatusOK, worldSum, fmt.Sprintf("%s/foo.zip", prefix), "foo.zip"},
// 		{http.StatusBadRequest, "", "", ""},
// 	}

// 	for _, tc := range testCases {
// 		bdy, err := json.Marshal(map[string]string{"key": tc.key})
// 		if err != nil {
// 			t.Fatal(err)
// 		}

// 		req := httptest.NewRequest(http.MethodPost, "/", bytes.NewBuffer(bdy))
// 		w := httptest.NewRecorder()

// 		handler.ServeHTTP(w, req)

// 		var d map[string]string
// 		resp := w.Result()

// 		if err := json.NewDecoder(resp.Body).Decode(&d); err != nil {
// 			t.Fatal(err)
// 		}

// 		if tc.wantCode != resp.StatusCode {
// 			t.Errorf("want status %d, got %d", tc.wantCode, resp.StatusCode)
// 		}

// 		gotSum := d["md5"]
// 		if tc.wantSum != gotSum {
// 			t.Errorf("want md5 sum '%s', got '%s'", tc.wantSum, gotSum)
// 		}

// 		gotLocation := d["location"]
// 		if tc.wantLocation != gotLocation {
// 			t.Errorf("want location '%s', got '%s'", tc.wantLocation, gotLocation)
// 		}
// 	}
// }
