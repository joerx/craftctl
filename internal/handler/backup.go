package handler

import (
	"bytes"
	"context"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io"
	"io/fs"
	"joerx/minecraft-cli/internal/zipper"
	"log"
	"net/http"
	"strings"
)

type ObjectStore interface {
	Put(ctx context.Context, key string, r io.Reader) (ObjectInfo, error)
}

type ObjectInfo struct {
	Location string
}

type Backup struct {
	rcon    RCon
	worldFS fs.FS
	store   ObjectStore
}

type backupRequest struct {
	Key string `json:"key"`
}

type backupResponse struct {
	MD5      string `json:"md5"`
	Location string `json:"location"`
}

func NewBackup(rc RCon, worldFS fs.FS, store ObjectStore) *Backup {
	return &Backup{rcon: rc, worldFS: worldFS, store: store}
}

// sanitizeKey replaces spaces in the filename and ensures the result has the desired suffix
func sanitizeKey(key, suffix string) string {
	key = strings.ReplaceAll(key, " ", "-")
	if !strings.HasSuffix(key, suffix) {
		key = fmt.Sprintf("%s%s", key, suffix)
	}
	return key
}

func (bh *Backup) zipAndStore(ctx context.Context, brq backupRequest) (backupResponse, error) {
	var resp backupResponse

	buf := new(bytes.Buffer)
	if err := zipper.ZipFS(bh.worldFS, buf); err != nil {
		return resp, err
	}

	oi, err := bh.store.Put(ctx, sanitizeKey(brq.Key, ".zip"), buf)
	if err != nil {
		return resp, err
	}

	checksum := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
	log.Printf("Archived world data, checksum is %s", checksum)

	resp.MD5 = checksum
	resp.Location = oi.Location

	return resp, nil
}

func (bh *Backup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Received backup request")

	var brq backupRequest
	if err := json.NewDecoder(r.Body).Decode(&brq); err != nil {
		serveJSONError(w, err, http.StatusInternalServerError)
		return
	}

	log.Printf("Storage key is '%s'", brq.Key)
	if brq.Key == "" {
		serveJSONError(w, fmt.Errorf("invalid storage key '%s'", brq.Key), http.StatusBadRequest)
		return
	}

	log.Println("Telling server to save the game")
	if err := bh.rcon.Command("save-all flush"); err != nil {
		serveJSONError(w, err, http.StatusInternalServerError)
		return
	}

	resp, err := bh.zipAndStore(r.Context(), brq)
	if err != nil {
		serveJSONError(w, err, http.StatusInternalServerError)
		return
	}

	serveJSON(w, resp, http.StatusOK)
}
