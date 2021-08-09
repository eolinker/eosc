package main

import (
	"bytes"
	"encoding/gob"
	"go.etcd.io/etcd/client/pkg/v3/types"
	"log"
)

type commandType int32

const INIT commandType = 1
const PROPOSE commandType = 2

// json接口交互用的结构

type Response struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
	Data []byte `json:"data"`
}

type ProposeMsg struct {
	From int `json:"from"`
	To int `json:"to"`
	Cmd  string `json:"cmd"`
	Data []byte  `json:"data"`
}

type JoinMsg struct {
	Id int `json:"id"`
	Host string `json:"host"`
	Peers map[types.ID]string `json:"peers"`
}


// Message 发送Propose和Init消息结构
type Message struct {
	Type	commandType
	From   int
	Cmd  string
	Data []byte
}

// node间进行通信的结构

// SnapStore 用于快照处理的结构
type SnapStore struct {
	Data        []byte
	Peer        map[types.ID]string
	ConfigChangeCount int
	Id          int
}


func (m *Message) Encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(m); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (m *Message) Decode(data []byte) error {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(m); err != nil {
		log.Fatalf("eosc: could not decode message (%v)", err)
		return err
	}
	return nil
}
func (s *SnapStore) Encode() ([]byte, error) {
	var buf bytes.Buffer
	if err := gob.NewEncoder(&buf).Encode(s); err != nil {
		return nil, err
	}
	return buf.Bytes(), nil
}
func (s *SnapStore) Decode(data []byte) error {
	dec := gob.NewDecoder(bytes.NewBuffer(data))
	if err := dec.Decode(s); err != nil {
		log.Fatalf("eosc: could not decode message (%v)", err)
		return err
	}
	return nil
}
