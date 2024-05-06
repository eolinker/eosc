package client

import (
	"context"
	"fmt"
	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/cmd/proto"
	"github.com/eolinker/eosc/process-admin/model"
)

func (i *imlClient) List(ctx context.Context, profession string) ([]model.Object, error) {
	err := i.conn.send(cmd.WorkerList, profession)
	if err != nil {
		return nil, err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return nil, err
	}
	if recv.Type() != proto.ArrayReply {
		return nil, fmt.Errorf("expect array but get %s", proto.ReplyTypeString(recv.Type()))
	}
	var mos []model.Object
	err = recv.Scan(&mos)
	if err != nil {
		return nil, err
	}
	return mos, nil
}

func (i *imlClient) Get(ctx context.Context, id string) (model.Object, error) {
	err := i.conn.send(cmd.WorkerGet, id)
	if err != nil {
		return nil, err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return nil, err
	}
	var mo model.Object
	err = recv.Scan(&mo)
	if err != nil {
		return nil, err
	}
	return mo, nil
}

func (i *imlClient) Set(ctx context.Context, id string, value any) error {
	err := i.conn.send(cmd.WorkerSet, id, value)
	if err != nil {
		return err
	}
	return i.conn.recvOk()
}

func (i *imlClient) Del(ctx context.Context, id string) error {
	err := i.conn.send(cmd.WorkerDel, id)
	if err != nil {
		return err
	}
	return i.conn.recvOk()
}

func (i *imlClient) MatchLabels(ctx context.Context, profession string, labels map[string]string) ([]model.Object, error) {
	if len(labels) == 0 {
		return i.List(ctx, profession)
	}
	err := i.conn.send(cmd.WorkerMatch, profession, labels)
	if err != nil {
		return nil, err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return nil, err
	}
	if recv.Type() != proto.ArrayReply {
		return nil, fmt.Errorf("expect array but get %s", proto.ReplyTypeString(recv.Type()))
	}
	var mos []model.Object
	err = recv.Scan(&mos)
	if err != nil {
		return nil, err
	}
	return mos, nil
}

//func (i *imlClient) Exists(ctx context.Context, id string) (bool, error) {
//	//TODO implement me
//	panic("implement me")
//}
