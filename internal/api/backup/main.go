package backup

import (
	"context"
	"io/fs"
	"joerx/minecraft-cli/internal/api/rcon"
)

type InputError string

func (ie InputError) Error() string {
	return string(ie)
}

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
	ObjectInfo
	MD5 string `json:"md5"`
}

type ListBackupOutput struct {
	Objects []ObjectInfo
}

type RestoreBackupInput struct{}

type RestoreBackupOutput struct {
}
