package client

import (
	"github.com/eolinker/eosc/process-admin/model"
)

type Client interface {
	Get(id string) (any, error)
	Set(id string, value any) error
	Del(id string) error
	Exists(id string) (bool, error)
	PList() ([]*model.ProfessionInfo, error)
	PGet(name string) (*model.ProfessionConfig, error)
	SGet(name string) (any, error)
	SSet(name string, value any) error
	VAll() (map[string]model.Variables, error)
	VGet(namespace string) (model.Variables, error)
	VSet(namespace string, v model.Variables) error
	Begin() error
}

type Transaction interface {
	Client
	Commit() error
	Rollback() error
}
