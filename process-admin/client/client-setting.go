package client

import (
	"context"
	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/model"
)

func (i *imlClient) SGet(ctx context.Context, name string) (model.Object, error) {
	err := i.conn.send(cmd.SGet, name)
	if err != nil {
		return nil, err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return nil, err
	}
	var r model.Object
	err = recv.Scan(&r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (i *imlClient) SSet(ctx context.Context, name string, value any) error {
	err := i.conn.send(cmd.SSet, name, value)
	if err != nil {
		return err
	}
	return i.conn.recvOk()
}
