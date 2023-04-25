package srv

import (
	"context"
	"joerx/minecraft-cli/frontend"
	"joerx/minecraft-cli/handler"
	"joerx/minecraft-cli/mc"
	"joerx/minecraft-cli/systemd"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Opts struct {
	RCONHostPort string
	RCONPasswd   string
	MCWorldDir   string
	Addr         string
	UnitName     string
}

func newServer(opts Opts) (*http.ServeMux, error) {
	log.Printf("%#v", opts)

	client, err := mc.NewClient(mc.ClientOpts{Password: opts.RCONPasswd, HostPort: opts.RCONHostPort})
	if err != nil {
		return nil, err
	}

	ctx := context.Background()

	uc, err := systemd.NewUnitController(ctx, opts.UnitName)
	if err != nil {
		return nil, err
	}

	worldFS := os.DirFS(opts.MCWorldDir)

	mux := http.NewServeMux()

	mux.Handle("/", frontend.New())
	mux.Handle("/cmd", handler.NewCommand(client))
	mux.Handle("/backup", handler.NewBackup(client, worldFS))
	mux.Handle("/start", handler.NewStart(uc))
	mux.Handle("/stop", handler.NewStop(uc))
	mux.Handle("/status", handler.NewStatus(uc))

	return mux, nil
}

func runServer(m http.Handler, addr string) error {
	if err := http.ListenAndServe(addr, m); err != nil {
		log.Fatal(err)
	}
	return nil
}

func Run(opts Opts) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	mux, err := newServer(opts)
	if err != nil {
		return err
	}

	go runServer(mux, opts.Addr)

	sig := <-sigs
	log.Printf("Received %v, shutting down", sig)

	return nil
}
