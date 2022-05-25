package raft

import (
	"bytes"
	"encoding/gob"
	"encoding/json"
	"log"

	"go.etcd.io/etcd/client/pkg/v3/types"
)


// json接口交互用的结构

type Response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
	Err  string      `json:"error,omitempty"`
}

type ProposeMsg struct {
	From uint64 `json:"from"`
	Body []byte `json:"body"`
}

type PostNodeRequest struct {
	BroadcastIP   string `json:"broadcast_ip"`
	BroadcastPort int    `json:"broadcast_port"`
	Protocol      string `json:"protocol"`
	Target        string `json:"target"`
}

type SNRequest struct {
	BroadcastIP   string `json:"broadcast_ip"`
	BroadcastPort int    `json:"broadcast_port"`
	Protocol      string `json:"protocol"`
	Target        string `json:"target"`
}

type SNResponse struct {
	SN string `json:"lastSN"`
}

type JoinRequest struct {
	NodeID        uint64 `json:"node_id"`
	NodeKey       string `json:"node_key"`
	BroadcastIP   string `json:"broadcast_ip"`
	BroadcastPort int    `json:"broadcast_port"`
	Protocol      string `json:"protocol"`
	Target        string `json:"target"`
}

type JoinResponse struct {
	*NodeSecret
	Peer map[uint64]*NodeInfo `json:"peer"`
}

type NodeSecret struct {
	ID  uint64 `json:"id"`
	Key string `json:"key"`
}

type NodeInfo struct {
	*NodeSecret
	BroadcastIP   string `json:"broadcast_ip"`
	BroadcastPort int    `json:"broadcast_port"`
	Addr          string `json:"addr"`
	Protocol      string `json:"protocol"`
}

func (n *NodeInfo) Marshal() []byte {
	data, _ := json.Marshal(n)
	return data
}


// node间进行通信的结构

// SnapStore 用于快照处理的结构
type SnapStore struct {
	Data              []byte
	Peer              map[types.ID]string
	ConfigChangeCount int
	Id                int
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
