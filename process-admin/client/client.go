package client

import (
	"context"
	"fmt"

	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/model"
)

type HClient interface {
	HGet(ctx context.Context, key string, field string) (string, error)
	HGetAll(ctx context.Context, key string) (map[string]string, error)
	HSet(ctx context.Context, key string, field string, value string) error
	HMSet(ctx context.Context, key string, fvs map[string]string) error
	HDelAll(ctx context.Context, key string) error
	HDel(ctx context.Context, key string, field string) error
	HExists(ctx context.Context, key string, field string) (bool, error)
	HKeys(ctx context.Context, key string) ([]string, error)
	HMatch(ctx context.Context, key string) ([]string, error)
}
type WorkerClient interface {
	List(ctx context.Context, profession string) ([]model.Object, error)
	Get(ctx context.Context, id string) (model.Object, error)
	Set(ctx context.Context, id string, value any) error
	Del(ctx context.Context, id string) error
	MatchLabels(ctx context.Context, profession string, labels map[string]string) ([]model.Object, error)
}
type ProfessionClient interface {
	PList(ctx context.Context) ([]*model.ProfessionInfo, error)
	PGet(ctx context.Context, name string) (*model.ProfessionConfig, error)
}
type SettingClient interface {
	SGet(ctx context.Context, name string) (model.Object, error)
	SSet(ctx context.Context, name string, value any) error
}
type VariableClient interface {
	VAll(ctx context.Context) (map[string]model.Variables, error)
	VGet(ctx context.Context, namespace string) (model.Variables, error)
	Value(ctx context.Context, namespace string, key string) (string, error)
	VSet(ctx context.Context, namespace string, v model.Variables) error
}
type Client interface {
	WorkerClient
	HClient
	ProfessionClient
	SettingClient
	VariableClient
	Ping(ctx context.Context) error

	//Exists(ctx context.Context, id string) (bool, error)

	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	Close() error
}
type imlClient struct {
	conn *baseClient
}

func (i *imlClient) Close() error {
	return i.conn.Close()
}

func (i *imlClient) Ping(ctx context.Context) error {
	err := i.conn.send(cmd.PING)
	if err != nil {
		return err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return err
	}
	return isPong(recv)
}

func (i *imlClient) Begin(ctx context.Context) error {
	err := i.conn.send(cmd.Begin)
	if err != nil {
		return err
	}
	return i.conn.recvOk()
}

func (i *imlClient) Commit(ctx context.Context) error {
	err := i.conn.send(cmd.Commit)
	if err != nil {
		return err
	}
	return i.conn.recvOk()
}

func (i *imlClient) Rollback(ctx context.Context) error {
	err := i.conn.send(cmd.Rollback)
	if err != nil {
		return err
	}
	return i.conn.recvOk()
}

func New(addrs ...string) (Client, error) {
	for _, addr := range addrs {
		c, err := create(addr)
		if err != nil {
			continue
		}
		return &imlClient{
			conn: c,
		}, nil
	}
	return nil, fmt.Errorf("no available address: [%v]", addrs)
}
