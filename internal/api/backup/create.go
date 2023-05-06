package backup

import (
	"bytes"
	"context"
	"crypto/md5"
	"fmt"
	"joerx/minecraft-cli/internal/zipper"
	"log"
	"os"
	"strings"
)

func (s *backupService) Create(ctx context.Context, in CreateBackupInput) (CreateBackupOutput, error) {
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

// sanitizeKey replaces spaces in the filename and ensures the result has the desired suffix
func sanitizeKey(key, suffix string) string {
	key = strings.ReplaceAll(key, " ", "-")
	if !strings.HasSuffix(key, suffix) {
		key = fmt.Sprintf("%s%s", key, suffix)
	}
	return key
}

func (s *backupService) zipAndStore(ctx context.Context, key string) (CreateBackupOutput, error) {
	buf := new(bytes.Buffer)
	worldFS := os.DirFS(s.worldDir)
	if err := zipper.ZipFS(worldFS, buf); err != nil {
		return CreateBackupOutput{}, err
	}

	checksum := fmt.Sprintf("%x", md5.Sum(buf.Bytes()))
	log.Printf("Archived world, checksum is %s", checksum)

	oi, err := s.store.Put(ctx, sanitizeKey(key, ".zip"), buf)
	if err != nil {
		return CreateBackupOutput{}, err
	}

	log.Printf("Uploaded world to %s", oi.Location)

	return CreateBackupOutput{
		MD5:        checksum,
		ObjectInfo: oi,
	}, nil
}
