package backup

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"io"
	"io/fs"
	"joerx/minecraft-cli/internal/service/rcon"
	"joerx/minecraft-cli/internal/zipper"
	"log"
	"strings"
)

type ObjectInfo struct {
	Location string
	Key      string
}

type Putter interface {
	Put(ctx context.Context, key string, r io.Reader) (ObjectInfo, error)
}

type Getter interface {
	Get(ctx context.Context, key string, w io.Writer) (ObjectInfo, error)
}

type Lister interface {
	List(ctx context.Context) ([]ObjectInfo, error)
}

type Store interface {
	Putter
	// Getter
	// Lister
}

type Config struct {
	RCon  rcon.RCon
	Store Store
	World fs.FS
}

func NewService(cfg Config) *Service {
	return &Service{cfg.RCon, cfg.Store, cfg.World}
}

type Service struct {
	rcon  rcon.RCon
	store Store
	world fs.FS
}

type CreateBackupInput struct {
	Key string `json:"key"`
}

type CreateBackupOutput struct {
	MD5      string `json:"md5"`
	Location string `json:"location"`
}

type ListBackupOutput struct {
	Objects []ObjectInfo
}

type RestoreBackupInput struct{}

type RestoreBackupOutput struct {
}

func (s *Service) Create(ctx context.Context, in CreateBackupInput) (CreateBackupOutput, error) {
	log.Println("Received backup request")
	log.Printf("Storage key is '%s'", in.Key)

	if in.Key == "" {
		return CreateBackupOutput{}, fmt.Errorf("invalid storage key '%s'", in.Key)
	}

	log.Println("Telling server to save the game")
	if _, err := s.rcon.Command(ctx, "save-all flush"); err != nil {
		return CreateBackupOutput{}, err
	}

	return s.zipAndStore(ctx, in.Key)
}

func (s *Service) List(ctx context.Context) (ListBackupOutput, error) {
	return ListBackupOutput{}, nil
}

func (s *Service) Restore(ctx context.Context, in RestoreBackupInput) (RestoreBackupOutput, error) {
	return RestoreBackupOutput{}, nil
}

// sanitizeKey replaces spaces in the filename and ensures the result has the desired suffix
func sanitizeKey(key, suffix string) string {
	key = strings.ReplaceAll(key, " ", "-")
	if !strings.HasSuffix(key, suffix) {
		key = fmt.Sprintf("%s%s", key, suffix)
	}
	return key
}

func (s *Service) zipAndStore(ctx context.Context, key string) (CreateBackupOutput, error) {
	buf := new(bytes.Buffer)
	if err := zipper.ZipFS(s.world, buf); err != nil {
		return CreateBackupOutput{}, err
	}

	oi, err := s.store.Put(ctx, sanitizeKey(key, ".zip"), buf)
	if err != nil {
		return CreateBackupOutput{}, err
	}

	checksum := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
	log.Printf("Archived world data, checksum is %s", checksum)

	return CreateBackupOutput{
		MD5:      checksum,
		Location: oi.Location,
	}, nil
}
