package admin_o

import (
	"context"
	open_api "github.com/eolinker/eosc/open-api"
	"github.com/eolinker/eosc/professions"
)

type AdminController interface {
	Transaction(context.Context, func(ctx context.Context, api AdminApi) error) ([]*open_api.EventResponse, error)
	Begin(ctx context.Context) (AdminTransaction, error)
}
type AdminTransaction interface {
	AdminApi
	Commit() ([]*open_api.EventResponse, error)
	Rollback() error
}
type AdminApi interface {
	ListWorker(ctx context.Context, profession string) ([]*WorkerInfo, error)
	GetWorker(ctx context.Context, id string) (*WorkerInfo, error)
	DeleteWorker(ctx context.Context, id string) (*WorkerInfo, error)
	SetWorker(ctx context.Context, profession, name, driver, version, desc string, data IData) (*WorkerInfo, error)
	AllWorkers(ctx context.Context) map[string][]*WorkerInfo

	GetProfession(ctx context.Context, profession string) (*professions.Profession, bool)
	ListProfession(ctx context.Context, profession string) ([]*professions.Profession, error)

	GetSetting(ctx context.Context, name string) (any, bool)
	SetSetting(ctx context.Context, name string, data IData) error

	AllVariables(ctx context.Context) map[string]map[string]string
	GetVariables(ctx context.Context, namespace string) (map[string]string, bool)
	GetVariable(ctx context.Context, namespace, key string) (string, bool)
	SetVariable(ctx context.Context, namespace string, values map[string]string) error
}
