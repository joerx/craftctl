package backup

import (
	"context"
	"io/fs"
	"joerx/minecraft-cli/internal/service/rcon"
)

type Config struct {
	RCon  rcon.RCon
	Store Store
	World fs.FS
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
	MD5      string `json:"md5"`
	Location string `json:"location"`
}

type ListBackupOutput struct {
	Objects []ObjectInfo
}

type RestoreBackupInput struct{}

type RestoreBackupOutput struct {
}
