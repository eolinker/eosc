package raft

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"strconv"
	"sync"
	"time"
)

// Demo Service
type service struct {
	store map[string]string
	mutex sync.RWMutex
}

// KV 用于传输的结构
type KV struct {
	Key   string
	Value string
}

func Create() *service {
	return &service{
		store: make(map[string]string),
		mutex: sync.RWMutex{},
	}
}

func (s *service) CommitHandler(cmd string, data []byte) (err error) {
	// TODO: process the command
	s.mutex.Lock()
	defer s.mutex.Unlock()
	switch cmd {
	case "set":
		kv := &KV{}
		err = json.Unmarshal(data, kv)
		s.store[kv.Key] = kv.Value
		return err
	case "init":
		return nil
	}
	return nil
}

func (s *service) ProcessHandler(command string, propose []byte) (cmd string, data []byte, err error) {
	// TODO: process the command before sending the message
	kv := &KV{}
	kv.Decode(propose)
	kv.Value = kv.Value + strconv.FormatInt(time.Now().UnixNano(), 10)
	data, err = json.Marshal(kv)
	command = "set"
	return command, data, err
}

func (s *service) GetInit() (cmd string, data []byte, err error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	cmd = "init"
	data, err = json.Marshal(s.store)
	return cmd, data, err
}

func (s *service) ResetSnap(data []byte) error {
	s.mutex.Lock()
	defer s.mutex.Unlock()
	store := make(map[string]string)
	json.Unmarshal(data, &store)
	s.store = store
	return nil
}

func (s *service) GetSnapshot() ([]byte, error) {
	s.mutex.RLock()
	defer s.mutex.RUnlock()
	return json.Marshal(s.store)
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
