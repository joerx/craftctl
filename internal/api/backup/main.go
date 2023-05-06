package backup

import (
	"context"
	"joerx/minecraft-cli/internal/api/rcon"
	"joerx/minecraft-cli/internal/systemd"
)

type InputError string

func (ie InputError) Error() string {
	return string(ie)
}

type Config struct {
	RCon     rcon.RCon
	Store    Store
	WorldDir string
	UC       *systemd.UnitController
}

type Service interface {
	Create(context.Context, CreateBackupInput) (CreateBackupOutput, error)
	List(context.Context) (ListBackupOutput, error)
	Restore(context.Context, RestoreBackupInput) (RestoreBackupOutput, error)
}

type CreateBackupInput struct {
	Key string `json:"key"`
}

type CreateBackupOutput struct {
	ObjectInfo
	MD5 string `json:"md5"`
}

type ListBackupOutput struct {
	Backups []ObjectInfo `json:"backups"`
}

type RestoreBackupInput struct {
	Key string `json:"key"`
}

type RestoreBackupOutput struct {
	Message string `json:"message"`
}

func NewService(cfg Config) Service {
	return &backupService{cfg.RCon, cfg.Store, cfg.UC, cfg.WorldDir}
}

type backupService struct {
	rcon     rcon.RCon
	store    Store
	uc       *systemd.UnitController
	worldDir string
}
