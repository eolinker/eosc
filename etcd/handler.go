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

func (s *_Server) addHandler(mux *http.ServeMux) {
	emptyHandler := http.NotFoundHandler()
	s.raftHandler.Store(&emptyHandler)
	s.leaseHandler.Store(&emptyHandler)
	s.hashKVHandler.Store(&emptyHandler)
	s.downgradeEnabledHandler.Store(&emptyHandler)

	mux.HandleFunc("/raft/node/join", s.join)

	mux.HandleFunc(rafthttp.RaftPrefix, func(w http.ResponseWriter, r *http.Request) {
		handler := s.raftHandler.Load()
		(*handler).ServeHTTP(w, r)
	})
	mux.HandleFunc(rafthttp.RaftPrefix+"/", func(w http.ResponseWriter, r *http.Request) {
		raftHandler := s.raftHandler.Load()
		(*raftHandler).ServeHTTP(w, r)
	})
	mux.HandleFunc(peerMembersPath, s.peerMembersHandler)
	mux.HandleFunc(peerMemberPromotePrefix, s.peerMemberPromoteHandler)

	leaseHandler := func(w http.ResponseWriter, r *http.Request) {
		leaseHandler := s.leaseHandler.Load()
		(*leaseHandler).ServeHTTP(w, r)
	}

	mux.HandleFunc(leasehttp.LeasePrefix, leaseHandler)
	mux.HandleFunc(leasehttp.LeaseInternalPrefix, leaseHandler)

	downgradeEnabledHandler := func(w http.ResponseWriter, r *http.Request) {

		downgradeEnabledHandler := s.downgradeEnabledHandler.Load()
		(*downgradeEnabledHandler).ServeHTTP(w, r)
	}

	mux.HandleFunc(etcdserver.DowngradeEnabledPath, downgradeEnabledHandler)

	mux.HandleFunc(etcdserver.PeerHashKVPath, func(w http.ResponseWriter, r *http.Request) {
		hashKVHandler := s.hashKVHandler.Load()
		(*hashKVHandler).ServeHTTP(w, r)

	})
	mux.HandleFunc("/version", func(w http.ResponseWriter, r *http.Request) {

		if !allowMethod(w, r, "GET") {
			return
		}
		vs := s.Version()

		w.Header().Set("Content-Type", "application/json")
		b, _ := json.Marshal(&vs)

		w.Write(b)

	})
}

type joinRequest struct {
	Addr   []string `json:"addr"`
	Name   string   `json:"name"`
	Client []string `json:"client"`
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
	response.Members, err = s.addMember(request.Name, request.Addr, request.Client)
	if err != nil {
		panic(err)
	}
	writeSuccessResult(w, response)
}

func (s *_Server) addMember(name string, urls []string, clients []string) (map[string][]string, error) {
	purls, err := types.NewURLs(urls)
	if err != nil {
		return nil, err
	}
	now := time.Now()
	member := membership.NewMember(name, purls, "APINTO_CLUSTER", &now)
	member.ClientURLs = clients
	ctx, _ := s.requestContext()
	members, err := s.server.AddMember(ctx, *member)

	res := make(map[string][]string)
	for _, m := range members {
		res[m.Name] = m.PeerURLs
	}

	//InitialCluster := initialClusterString(res)
	//
	//s.resetCluster(InitialCluster)
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
