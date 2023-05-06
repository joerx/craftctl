package handler

import (
	"context"
	"encoding/json"
	"joerx/minecraft-cli/internal/api/backup"
	"net/http"
)

// Using a closure is a simple way to "inject" dependencies into our handler
// We avoid repetitive error handling code by using our custom http.Handler type
// See https://go.dev/blog/error-handling-and-go
func CreateBackup(svc func(context.Context, backup.CreateBackupInput) (backup.CreateBackupOutput, error)) http.Handler {
	return h(func(w http.ResponseWriter, r *http.Request) error {
		var input backup.CreateBackupInput
		if err := json.NewDecoder(r.Body).Decode(&input); err != nil {
			return err
		}

		output, err := svc(r.Context(), input)
		if _, ok := err.(backup.InputError); ok {
			return badRequest(err)
		}
		if err != nil {
			return err
		}

		serveJSON(w, output, http.StatusOK)
		return nil
	})
}
