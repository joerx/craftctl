package server

import (
	"joerx/minecraft-cli/frontend"
	"joerx/minecraft-cli/internal/handler"
	"net/http"
)

func newRouter(app *application) *http.ServeMux {
	mux := http.NewServeMux()

	mux.Handle("/", frontend.New())
	mux.Handle("/cmd", handler.Command(app.RCon.Command))
	mux.Handle("/backup/create", handler.CreateBackup(app.Backup.Create))
	mux.Handle("/backup/list", handler.ListBackups(app.Backup.List))
	mux.Handle("/backup/restore", handler.RestoreBackup(app.Backup.Restore))
	mux.Handle("/start", handler.NewStart(app.UC))
	mux.Handle("/stop", handler.NewStop(app.UC))
	mux.Handle("/status", handler.NewStatus(app.UC))

	return mux
}
