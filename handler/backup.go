package handler

import (
	"bytes"
	"crypto/md5"
	"encoding/json"
	"fmt"
	"io/fs"
	"joerx/minecraft-cli/zipper"
	"log"
	"net/http"
)

type Backup struct {
	rcon      RCon
	worldFS   fs.FS
	backupDir string
}

type backupResponse struct {
	MD5 string `json:"md5"`
}

func NewBackup(rc RCon, worldFS fs.FS) *Backup {
	return &Backup{rcon: rc, worldFS: worldFS}
}

func (bh *Backup) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	log.Println("Received backup request")

	log.Println("Telling server to save the game")
	if err := bh.rcon.Command("save-all flush"); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	buf := new(bytes.Buffer)
	if err := zipper.ZipFS(bh.worldFS, buf); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
	}

	checksum := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
	log.Printf("Archived world data, checksum is %s", checksum)

	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-type", "application/json")

	json.NewEncoder(w).Encode(backupResponse{MD5: checksum})
}
