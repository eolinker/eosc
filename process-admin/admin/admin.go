package admin

import (
	"context"
	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/process-admin/marshal"
	"github.com/eolinker/eosc/professions"
)

type AdminController interface {
	IAdminApiReader

	Transaction(context.Context, func(ctx context.Context, api AdminApiWrite) error) error
	Begin(ctx context.Context) AdminTransaction
}
type IAdminApiReader interface {
	eosc.IWorkers
	ListWorker(ctx context.Context, profession string) ([]*WorkerInfo, error)
	GetWorker(ctx context.Context, id string) (*WorkerInfo, error)
	AllWorkers(ctx context.Context) []*WorkerInfo

	GetProfession(ctx context.Context, profession string) (*professions.Profession, bool)
	ListProfession(ctx context.Context) []*professions.Profession

	GetSetting(ctx context.Context, name string) (any, bool)

	AllVariables(ctx context.Context) map[string]map[string]string
	GetVariables(ctx context.Context, namespace string) (map[string]string, bool)
	GetVariable(ctx context.Context, namespace, key string) (string, bool)
}
type AdminTransaction interface {
	AdminApiWrite
	IAdminApiReader
	Commit() error
	Rollback() error
}
type AdminApiWrite interface {
	IAdminApiReader

	SetProfession(name string, profession *eosc.ProfessionConfig) error
	ResetProfession(configs []*eosc.ProfessionConfig)
	DeleteWorker(ctx context.Context, id string) (*WorkerInfo, error)
	SetWorker(ctx context.Context, profession, name, driver, version, desc string, data marshal.IData) (*WorkerInfo, error)

	SetSetting(ctx context.Context, name string, data marshal.IData) error
	SetVariable(ctx context.Context, namespace string, values map[string]string) error
}
