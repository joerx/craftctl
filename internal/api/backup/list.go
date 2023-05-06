package backup

import "context"

// List returns a list of backups from the underlying store.
func (s *backupService) List(ctx context.Context) (ListBackupOutput, error) {
	backups, err := s.store.List(ctx)
	if err != nil {
		return ListBackupOutput{}, err
	}
	return ListBackupOutput{Backups: backups}, nil
}
