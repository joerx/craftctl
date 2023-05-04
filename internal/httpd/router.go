package httpd

import (
	"joerx/minecraft-cli/frontend"
	"joerx/minecraft-cli/internal/handler"
	"net/http"
)

func newRouter(app *Application) *http.ServeMux {
	mux := http.NewServeMux()

	commandHandler := handler.NewCommandHandler(app.RCon)
	backupHandler := handler.NewBackupHandler(app.Backup)

	mux.Handle("/", frontend.New())
	mux.Handle("/cmd", handler.Handler(commandHandler.RunCommand))
	mux.Handle("/backup", handler.Handler(backupHandler.Create))
	mux.Handle("/start", handler.NewStart(app.UC))
	mux.Handle("/stop", handler.NewStop(app.UC))
	mux.Handle("/status", handler.NewStatus(app.UC))

	return mux
}
