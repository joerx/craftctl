package server

import (
	"context"
	"fmt"
	"joerx/minecraft-cli/internal/api/backup"
	"joerx/minecraft-cli/internal/api/rcon"
	"joerx/minecraft-cli/internal/mc"
	"joerx/minecraft-cli/internal/storage/s3"
	"joerx/minecraft-cli/internal/systemd"
	"os"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
)

type appConfig struct {
	RCONHostPort string
	RCONPasswd   string
	MCWorldDir   string
	UnitName     string
	S3Bucket     string
	S3Region     string
}

type application struct {
	RCon   *rcon.Service
	Backup backup.Service
	UC     *systemd.UnitController
}

func newApp(cfg appConfig) (*application, error) {
	rc := mc.NewClient(mc.ClientConfig{Password: cfg.RCONPasswd, HostPort: cfg.RCONHostPort})

	uc, err := systemd.NewUnitController(context.Background(), cfg.UnitName)
	if err != nil {
		return nil, err
	}

	worldFS := os.DirFS(cfg.MCWorldDir)
	store, err := newStore(cfg)
	if err != nil {
		return nil, err
	}

	rsvc := rcon.NewService(rc)
	bsvc := backup.NewService(backup.Config{RCon: rc, World: worldFS, Store: store})

	return &application{
		RCon:   rsvc,
		Backup: bsvc,
		UC:     uc,
	}, nil
}

func newStore(cfg appConfig) (backup.Store, error) {
	if cfg.S3Region == "" {
		return nil, fmt.Errorf("no s3 region provided")
	}
	if cfg.S3Bucket == "" {
		return nil, fmt.Errorf("no s3 bucket provided")
	}

	sess, err := session.NewSession(&aws.Config{
		Region: &cfg.S3Region,
	})

	if err != nil {
		return nil, err
	}
	return s3.NewStore(sess, cfg.S3Bucket), nil
}
