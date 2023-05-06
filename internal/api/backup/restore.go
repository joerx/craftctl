package backup

import (
	"context"
	"fmt"
	"joerx/minecraft-cli/internal/zipper"
	"log"
	"os"
)

func (s *backupService) Restore(ctx context.Context, in RestoreBackupInput) (RestoreBackupOutput, error) {
	var out RestoreBackupOutput

	// Download backup to tmp
	tf, err := os.CreateTemp("", "world-")
	if err != nil {
		return out, err
	}

	if err := s.store.Get(ctx, in.Key, tf); err != nil {
		return out, err
	}

	log.Printf("Fetched %s to %s", in.Key, tf.Name())

	// Stop server
	log.Println("Stopping server")
	if _, err := s.uc.Stop(ctx); err != nil {
		return out, err
	}

	// Delete existing world directory
	if err := rmdir(s.worldDir); err != nil {
		return out, err
	}

	// Extract archive
	log.Printf("Extracting %s to %s", tf.Name(), s.worldDir)
	if err := zipper.Unzip(tf.Name(), s.worldDir); err != nil {
		return out, err
	}

	// Start server
	log.Println("Starting server")
	if _, err := s.uc.Start(ctx); err != nil {
		return out, err
	}

	msg := fmt.Sprintf("Restored world from backup %s", in.Key)
	log.Println(msg)

	out.Message = msg
	return out, nil
}

func rmdir(path string) error {
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		log.Printf("Directory %s exists, removing", path)
		if err := os.RemoveAll(path); err != nil {
			return err
		}
	}
	return nil
}
