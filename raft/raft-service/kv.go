package raft_service

import (
	"bytes"
	"encoding/gob"
)

// KV 用于传输的结构
type KV struct {
	Key    string
	Config interface{}
}

func (kv *KV) Encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(kv); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}

func (kv *KV) Decode(data []byte) error {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(kv); err != nil {
		return err
	}
	return nil
}
