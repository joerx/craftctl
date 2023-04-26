package handler

import (
	"bytes"
	"crypto/md5"
	"fmt"
	"io/fs"
	"joerx/minecraft-cli/internal/zipper"
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
		serveJSONError(w, err, http.StatusInternalServerError)
		return
	}

	buf := new(bytes.Buffer)
	if err := zipper.ZipFS(bh.worldFS, buf); err != nil {
		serveJSONError(w, err, http.StatusInternalServerError)
		return
	}

	checksum := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
	log.Printf("Archived world data, checksum is %s", checksum)

	serveJSON(w, backupResponse{MD5: checksum}, http.StatusOK)
}
