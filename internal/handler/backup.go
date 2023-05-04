package handler

import (
	"encoding/json"
	"joerx/minecraft-cli/internal/service/backup"
	"net/http"
)

type BackupHandler struct {
	svc *backup.Service
}

func NewBackupHandler(svc *backup.Service) *BackupHandler {
	return &BackupHandler{svc}
}

func (bh *BackupHandler) Create(w http.ResponseWriter, r *http.Request) error {
	var input backup.CreateBackupInput
	if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
		return err
	}

	output, err := bh.svc.Create(r.Context(), input)
	if err != nil {
		return err
	}

	serveJSON(w, output, http.StatusOK)
	return nil
}
