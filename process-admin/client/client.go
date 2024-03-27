package client

import (
	"context"
	"github.com/eolinker/eosc/process-admin/model"
)

type Client interface {
	List(ctx context.Context, profession string) ([]model.Object, error)
	Get(ctx context.Context, id string) (model.Object, error)
	Set(ctx context.Context, id string, value any) error
	Del(ctx context.Context, id string) error
	Exists(ctx context.Context, id string) (bool, error)
	PList(ctx context.Context) ([]*model.ProfessionInfo, error)
	PGet(ctx context.Context, name string) (*model.ProfessionConfig, error)
	SGet(ctx context.Context, name string) (any, error)
	SSet(ctx context.Context, name string, value any) error
	VAll(ctx context.Context) (map[string]model.Variables, error)
	VGet(ctx context.Context, namespace string) (model.Variables, error)
	VSet(ctx context.Context, namespace string, v model.Variables) error
	Transaction(ctx context.Context) (Transaction, error)
}

func New() Client {
	return nil
}

type Transaction interface {
	Client
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
}
