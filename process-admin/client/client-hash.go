package client

import (
	"context"
	"github.com/eolinker/eosc/process-admin/cmd"
)

func (i *imlClient) HGet(ctx context.Context, key string, field string) (string, error) {
	err := i.conn.send(cmd.HGet, key, field)
	if err != nil {
		return "", err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return "", err
	}
	return recv.String()
}

func (i *imlClient) HGetAll(ctx context.Context, key string) (map[string]string, error) {
	err := i.conn.send(cmd.HGetAll, key)
	if err != nil {
		return nil, err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return nil, err
	}
	rs, err := recv.Array()
	if err != nil {
		return nil, err
	}

	result := make(map[string]string)
	for len(rs) >= 2 {
		r := rs[:2]
		rs = rs[2:]
		var k, v string
		e := r.Scan(&k, &v)
		if e != nil {
			return nil, err
		}
		result[k] = v
	}
	return result, nil

}

func (i *imlClient) HSet(ctx context.Context, key string, field string, value string) error {
	err := i.conn.send(cmd.HSet, key, field, value)
	if err != nil {
		return err
	}
	return i.conn.recvOk()
}
func (i *imlClient) HRSet(ctx context.Context, key string, fvs map[string]string) error {
	args := make([]any, 0, len(fvs)*2+1)
	args = append(args, key)
	for k, v := range fvs {
		args = append(args, k, v)
	}
	err := i.conn.send(cmd.HRest, args...)
	if err != nil {
		return err
	}
	return i.conn.recvOk()
}
func (i *imlClient) HMSet(ctx context.Context, key string, fvs map[string]string) error {
	args := make([]any, 0, len(fvs)*2+1)
	args = append(args, key)
	for k, v := range fvs {
		args = append(args, k, v)
	}
	err := i.conn.send(cmd.HMSet, args...)
	if err != nil {
		return err
	}
	return i.conn.recvOk()
}

func (i *imlClient) HDelAll(ctx context.Context, key string) error {

	err := i.conn.send(cmd.HDelAll, key)
	if err != nil {
		return err
	}
	return i.conn.recvOk()
}

func (i *imlClient) HDel(ctx context.Context, key string, field string) error {
	err := i.conn.send(cmd.HDel, key, field)
	if err != nil {
		return err
	}
	return i.conn.recvOk()
}

func (i *imlClient) HExists(ctx context.Context, key string, field string) (bool, error) {
	err := i.conn.send(cmd.HExists, key, field)
	if err != nil {
		return false, err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return false, err
	}
	return recv.Bool()
}

func (i *imlClient) HKeys(ctx context.Context, key string) ([]string, error) {
	err := i.conn.send(cmd.HKeys, key)
	if err != nil {
		return nil, err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return nil, err
	}
	var keys []string
	err = recv.Scan(&keys)
	if err != nil {
		return nil, err
	}
	return keys, nil
}

func (i *imlClient) HMatch(ctx context.Context, key string) ([]string, error) {
	err := i.conn.send(cmd.HMatch, key)
	if err != nil {
		return nil, err
	}
	recv, err := i.conn.recv()
	if err != nil {
		return nil, err
	}
	var matches []string
	err = recv.Scan(&matches)
	if err != nil {
		return nil, err
	}
	return matches, nil
}
