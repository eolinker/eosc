package etcd

import (
	"encoding/json"
	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"go.etcd.io/etcd/server/v3/etcdserver/api/membership"
	"go.etcd.io/etcd/server/v3/etcdserver/api/rafthttp"
	"go.etcd.io/etcd/server/v3/lease/leasehttp"
	"io/ioutil"
	"net/http"
	"time"
)

const (
	peerMembersPath         = "/members"
	peerMemberPromotePrefix = "/members/promote/"
)

func (s *_Server) addHandler(mux *http.ServeMux)  {



	mux.HandleFunc("/raft/node/join", s.join)

	mux.HandleFunc(rafthttp.RaftPrefix, func(w http.ResponseWriter, r *http.Request) {
		s.mu.RLock()
		defer s.mu.RUnlock()
		if s.raftHandler != nil{
			s.raftHandler.ServeHTTP(w,r)
		}
	})
	mux.HandleFunc(rafthttp.RaftPrefix+"/", func(w http.ResponseWriter, r *http.Request) {
		s.mu.RLock()
		defer s.mu.RUnlock()
		if s.raftHandler != nil{
			s.raftHandler.ServeHTTP(w,r)
		}
	})
	mux.HandleFunc(peerMembersPath, s.peerMembersHandler)
	mux.HandleFunc(peerMemberPromotePrefix, s.peerMemberPromoteHandler)

	 leaseHandler:= func(w http.ResponseWriter, r *http.Request) {
		 s.mu.RLock()
		 defer s.mu.RUnlock()
		if s.leaseHandler!= nil{
			 s.leaseHandler.ServeHTTP(w,r)
			return
		}
		http.NotFound(w,r)
	 }

	 mux.HandleFunc(leasehttp.LeasePrefix, leaseHandler)
	 mux.HandleFunc(leasehttp.LeaseInternalPrefix, leaseHandler)

	downgradeEnabledHandler:= func(w http.ResponseWriter, r *http.Request) {
		s.mu.RLock()
		defer s.mu.RUnlock()
		if s.downgradeEnabledHandler!= nil{
			s.downgradeEnabledHandler.ServeHTTP(w,r)
			return
		}
		http.NotFound(w,r)
	}

	mux.HandleFunc(etcdserver.DowngradeEnabledPath, downgradeEnabledHandler)

	hashKVHandler:= s.server.HashKVHandler()
	if hashKVHandler != nil {
		mux.Handle(etcdserver.PeerHashKVPath, hashKVHandler)
	}
	mux.HandleFunc("/version", versionHandler(s.server.Cluster(), serveVersion))
}


type joinRequest struct {
	Addr []string `json:"addr"`
	Name string `json:"name"`
}
type joinResponse struct {
	Members map[string][]string `json:"members"`
}

type Response struct {
	Code string      `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data,omitempty"`
}
func (s *_Server) join(w http.ResponseWriter, r *http.Request) {
	data, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer func() {
		_ = r.Body.Close()
	}()
	request := new(joinRequest)
	err = json.Unmarshal(data, request)
	if err != nil {
		writeError(w, "110000", err.Error())
		return
	}
	response := new(joinResponse)
	response.Members,  err = s.addMember(request.Name,request.Addr)
	if err != nil {
		panic(err)
	}
	writeSuccessResult(w, response)
}

func (s *_Server) addMember(name string,urls []string) (map[string][]string,  error) {
	purls, err := types.NewURLs(urls)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	member := membership.NewMember("", purls, "apinto", &now)
	member.Name = name
	ctx, _ := s.requestContext()
	members, err := s.server.AddMember(ctx, *member)

	res := make(map[string][]string)
	for _, m := range members {
		res[m.Name] = m.PeerURLs
	}
	return res, nil
}

// writeSuccessResult 返回成功结果
func writeSuccessResult(w http.ResponseWriter, value interface{}) {
	result := &Response{
		Code: "000000",
		Msg:  "success",
		Data: value,
	}
	writeTo(w, result)
}

// writeError 返回失败结果
func writeError(w http.ResponseWriter, code string, errInfo string) {
	result := &Response{
		Code: code,
		Msg:  errInfo,
	}
	writeTo(w, result)
}

func writeTo(w http.ResponseWriter, obj interface{}) {
	if data, ok := obj.([]byte); ok {
		_, _ = w.Write(data)
		return
	}
	data, _ := json.Marshal(obj)
	_, _ = w.Write(data)
}

func decodeResponse(data []byte) (*Response, error) {
	res := new(Response)
	err := json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
