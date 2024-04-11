package client

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/eolinker/eosc/process-admin/cmd"
	"github.com/eolinker/eosc/process-admin/cmd/proto"
	"github.com/eolinker/eosc/process-admin/model"
)

type Client interface {
	Ping(ctx context.Context) error
	List(ctx context.Context, profession string) ([]model.Object, error)
	Get(ctx context.Context, id string) (model.Object, error)
	Set(ctx context.Context, id string, value any) error
	Del(ctx context.Context, id string) error
	//Exists(ctx context.Context, id string) (bool, error)
	PList(ctx context.Context) ([]*model.ProfessionInfo, error)
	PGet(ctx context.Context, name string) (*model.ProfessionConfig, error)
	SGet(ctx context.Context, name string) (model.Object, error)
	SSet(ctx context.Context, name string, value any) error
	VAll(ctx context.Context) (map[string]model.Variables, error)
	VGet(ctx context.Context, namespace string) (model.Variables, error)
	Value(ctx context.Context, namespace string, key string) (string, error)
	VSet(ctx context.Context, namespace string, v model.Variables) error
	Begin(ctx context.Context) error
	Commit(ctx context.Context) error
	Rollback(ctx context.Context) error
	MatchLabels(ctx context.Context, profession string, labels map[string]string) ([]model.Object, error)
}
type imlClient struct {
	conn *baseClient
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

//func (i *imlClient) Exists(ctx context.Context, id string) (bool, error) {
//	//TODO implement me
//	panic("implement me")
//}

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
