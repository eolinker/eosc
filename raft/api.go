package raft

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/julienschmidt/httprouter"
)

// 客户端请求处理
type Client struct {
	raft *Node
}

type jsonResponse struct {
	Code   string      `json:"code"`
	Msg    string      `json:"msg"`
	Result interface{} `json:"result"`
}

func (c *Client) Handler() http.Handler {
	router := httprouter.New()
	router.HandlerFunc("POST", "/raft/api/set", c.setHandler)
	router.HandlerFunc("POST", "/raft/api/deleteNode", c.deleteNodeHandler)
	router.HandlerFunc("POST", "/raft/api/addNode", c.addNodeHandler)
	router.HandlerFunc("GET", "/raft/api/getPeerList", c.getPeersHandler)
	router.HandlerFunc("GET", "/raft/api/get", c.getHandler)
	return router
}

func (c *Client) getPeersHandler(w http.ResponseWriter, r *http.Request) {
	res := &jsonResponse{
		Code: "000000",
		Msg:  "success",
	}
	list, count, err := c.raft.GetPeers()
	if err != nil {
		res.Code = "000001"
		res.Msg = err.Error()
	} else {
		res.Result = list
		res.Msg = strconv.Itoa(count)
	}
	c.writeResult(w, res)
}
func (c *Client) getHandler(w http.ResponseWriter, r *http.Request) {
	res := &jsonResponse{
		Code: "000000",
		Msg:  "success",
	}
	key := r.URL.Query().Get("key")
	if key == "" {
		res.Msg = "parameters format error"
		res.Code = "000001"
	} else {
		v, ok := c.raft.Service.(*service)
		if ok {
			val, ok := v.store[key]
			if ok {
				res.Msg = val
			} else {
				res.Msg = ""
			}
		} else {
			res.Msg = ""
		}
	}
	c.writeResult(w, res)
}

func (c *Client) setHandler(w http.ResponseWriter, r *http.Request) {
	res := &jsonResponse{
		Code: "000000",
		Msg:  "success",
	}
	key := r.PostFormValue("key")
	value := r.PostFormValue("value")
	if key == "" || value == "" {
		res.Msg = "parameters format error"
		res.Code = "111111"
	} else {
		kv := &KV{
			Key:   key,
			Value: value,
		}
		data, err := kv.Encode()
		if err != nil {
			res.Msg = err.Error()
			res.Code = "000001"
		}
		err = c.raft.Send("set", data)
		if err != nil {
			res.Msg = err.Error()
			res.Code = "000001"
		}
	}
	c.writeResult(w, res)
}

// deleteNodeHandler 删除节点
func (c *Client) deleteNodeHandler(w http.ResponseWriter, r *http.Request) {
	res := &jsonResponse{
		Code: "000000",
		Msg:  "success",
	}
	nodeId := r.PostFormValue("Id")
	if nodeId == "" {
		res.Msg = "parameters format error"
		res.Code = "000001"
	} else {
		Id, err := strconv.ParseUint(nodeId, 0, 64)
		if err != nil {
			res.Msg = fmt.Sprintf("Failed to convert ID for conf change (%v)\n", err)
			res.Code = "000002"
		} else {
			err = c.raft.DeleteConfigChange(Id)
			if err != nil {
				res.Code = "000003"
				res.Msg = err.Error()
			}
		}
	}
	c.writeResult(w, res)
}

// addNodeHandler 将节点加入集群
func (c *Client) addNodeHandler(w http.ResponseWriter, r *http.Request) {
	res := &jsonResponse{
		Code: "000000",
		Msg:  "success",
	}
	nodeId := r.PostFormValue("Id")
	addr := r.PostFormValue("host")
	if nodeId == "" {
		res.Msg = "parameters format error"
		res.Code = "000001"
	} else {
		// 获取节点ID
		Id, err := strconv.ParseUint(nodeId, 0, 64)
		if err != nil {
			res.Msg = fmt.Sprintf("Failed to convert ID for conf change (%v)\n", err)
			res.Code = "000002"
		} else {
			err = c.raft.AddConfigChange(Id, addr)
			if err != nil {
				res.Code = "000003"
				res.Msg = err.Error()
			}
		}
	}
	c.writeResult(w, res)
}

func (c *Client) writeResult(w http.ResponseWriter, v interface{}) {
	data, err := json.Marshal(v)
	if err != nil {
		w.WriteHeader(500)
		w.Write([]byte(err.Error()))
		return
	}
	w.Write(data)
}
