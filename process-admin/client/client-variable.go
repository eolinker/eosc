package client

import (
	"context"
	"encoding/json"
	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/model"
)

func (i *imlClient) VAll(ctx context.Context) (map[string]model.Variables, error) {
	err := i.conn.send(cmd.VGet)
	if err != nil {
		return nil, err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return nil, err
	}
	var data []byte
	err = recv.Scan(&data)
	if err != nil {
		return nil, err
	}
	var r map[string]model.Variables
	err = json.Unmarshal(data, &r)
	if err != nil {
		return nil, err
	}
	return r, nil
}

func (i *imlClient) VGet(ctx context.Context, namespace string) (model.Variables, error) {
	err := i.conn.send(cmd.VGet, namespace)
	if err != nil {
		return nil, err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return nil, err
	}

	var r model.Variables
	err = recv.Scan(&r)
	if err != nil {
		return nil, err
	}

	return r, nil
}
func (i *imlClient) Value(ctx context.Context, namespace string, key string) (string, error) {
	err := i.conn.send(cmd.VGet, namespace, key)
	if err != nil {
		return "", err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return "", err
	}
	return recv.String()

}
func (i *imlClient) VSet(ctx context.Context, namespace string, v model.Variables) error {
	err := i.conn.send(cmd.VSet, namespace, v)
	if err != nil {
		return err
	}
	return i.conn.recvOk()
}
