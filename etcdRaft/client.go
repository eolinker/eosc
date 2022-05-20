package etcdRaft

import (
	"context"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"log"
)

type IEtcdKVClient interface {
	Get(key string) ([]byte, error)
	GetPrefix(prefix string) (map[string][]byte, error)
	Put(key, value string) error
	Delete(key string) error
}

func (e *EtcdServer) requestContext() (context.Context, context.CancelFunc) {
	return context.WithTimeout(e.ctx, e.requestTimeout)
}

// Get 获取数据
func (e *EtcdServer) Get(key string) ([]byte, error) {
	kv, err := e.GetRaw(key)
	if err != nil || kv == nil {
		return nil, err
	}
	return kv.Value, nil
}

// GetRaw 获取Raw数据
func (e *EtcdServer) GetRaw(key string) (*mvccpb.KeyValue, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	client := e.client
	ctx, cancel := e.requestContext()
	defer cancel()
	resp, err := client.Get(ctx, key)
	if err != nil {
		return nil, err
	}
	if len(resp.Kvs) == 0 {
		return nil, nil
	}
	return resp.Kvs[0], nil
}

// GetPrefix 获取指定前缀的数据
func (e *EtcdServer) GetPrefix(prefix string) (map[string][]byte, error) {
	rawKVs, err := e.GetRawPrefix(prefix)
	if err != nil {
		return nil, err
	}
	kvs := make(map[string][]byte)
	for _, kv := range rawKVs {
		kvs[string(kv.Key)] = kv.Value
	}

	return kvs, nil
}

// GetRawPrefix 获取指定前缀的Raw数据
func (e *EtcdServer) GetRawPrefix(prefix string) (map[string]*mvccpb.KeyValue, error) {
	e.mu.RLock()
	defer e.mu.RUnlock()
	client := e.client
	resp, err := func() (*clientv3.GetResponse, error) {
		ctx, cancel := e.requestContext()
		defer cancel()
		return client.Get(ctx, prefix, clientv3.WithPrefix())
	}()
	if err != nil {
		return nil, err
	}
	kvs := make(map[string]*mvccpb.KeyValue)
	for _, kv := range resp.Kvs {
		kvs[string(kv.Key)] = kv
	}
	return kvs, nil
}

func (e *EtcdServer) Delete(key string) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	client := e.client
	ctx, cancel := e.requestContext()
	defer cancel()
	_, err := client.Delete(ctx, key)
	return err
}

// Put Put数据
func (e *EtcdServer) Put(key, value string) error {
	e.mu.RLock()
	defer e.mu.RUnlock()
	client := e.client

	ctx, cancel := e.requestContext()
	defer cancel()

	_, err := client.Put(ctx, key, value)

	return err
}

func (e *EtcdServer) getAllData() (map[string][]byte, error) {
	client := e.client
	resp, err := func() (*clientv3.GetResponse, error) {
		ctx, cancel := e.requestContext()
		defer cancel()
		return client.Get(ctx, dataPrefixKey, clientv3.WithPrefix())
	}()
	if err != nil {
		return nil, err
	}
	kvs := make(map[string][]byte)
	for _, kv := range resp.Kvs {
		kvs[string(kv.Key)] = kv.Value
	}
	return kvs, nil
}

func (e *EtcdServer) resetAllData(data map[string][]byte) {
	client := e.client
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	for key, bytes := range data {
		_, err := client.Put(ctx, key, string(bytes))
		if err != nil {
			log.Printf("reset all data error : %s", err.Error())
		}
	}
	return
}
