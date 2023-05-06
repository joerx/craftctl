package server

import (
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
)

type Config struct {
	appConfig
	Addr string
}

func newServer(cfg appConfig) (*http.ServeMux, error) {
	app, err := newApp(cfg)
	if err != nil {
		return nil, err
	}
	return newRouter(app), nil
}

func runServer(m http.Handler, addr string) error {
	if err := http.ListenAndServe(addr, m); err != nil {
		log.Fatal(err)
	}
	return nil
}

func Run(cfg Config) error {
	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGTERM, syscall.SIGINT)

	mux, err := newServer(cfg.appConfig)
	if err != nil {
		return err
	}

	go runServer(mux, cfg.Addr)

	sig := <-sigs
	log.Printf("Received %v, shutting down", sig)

	return nil
}
