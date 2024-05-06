package client

import (
	"context"
	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/model"
)

func (i *imlClient) PList(ctx context.Context) ([]*model.ProfessionInfo, error) {
	err := i.conn.send(cmd.PList)
	if err != nil {
		return nil, err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return nil, err
	}

	var rs []*model.ProfessionInfo
	err = recv.Scan(&rs)
	if err != nil {
		return nil, err
	}
	return rs, nil

}

func (i *imlClient) PGet(ctx context.Context, name string) (*model.ProfessionConfig, error) {
	err := i.conn.send(cmd.PGet, name)
	if err != nil {
		return nil, err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return nil, err
	}
	r := new(model.ProfessionConfig)
	err = recv.Scan(r)
	if err != nil {
		return nil, err
	}
	return r, nil
}
