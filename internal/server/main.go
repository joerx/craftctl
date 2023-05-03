package server

import (
	"context"
	"fmt"
	"joerx/minecraft-cli/internal/frontend"
	"joerx/minecraft-cli/internal/handler"
	"joerx/minecraft-cli/internal/mc"
	"joerx/minecraft-cli/internal/storage/s3"
	"joerx/minecraft-cli/internal/systemd"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type Opts struct {
	RCONHostPort string
	RCONPasswd   string
	MCWorldDir   string
	Addr         string
	UnitName     string
	S3Bucket     string
	S3Region     string
}

func newStore(opts Opts) (handler.ObjectStore, error) {
	if opts.S3Region == "" {
		return nil, fmt.Errorf("no s3 region provided")
	}
	if opts.S3Bucket == "" {
		return nil, fmt.Errorf("no s3 bucket provided")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: &opts.S3Region,
	})

	if err != nil {
		return nil, err
	}
	return s3.NewStore(sess, opts.S3Bucket), nil
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

	store, err := newStore(opts)
	if err != nil {
		return nil, err
	}

	worldFS := os.DirFS(opts.MCWorldDir)

	mux := http.NewServeMux()

	mux.Handle("/", frontend.New())
	mux.Handle("/cmd", handler.NewCommand(client))
	mux.Handle("/backup", handler.NewBackup(client, worldFS, store))
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
