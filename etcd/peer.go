package etcd

import (
	"encoding/json"
	"errors"
	"fmt"
	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/server/v3/etcdserver"
	"go.etcd.io/etcd/server/v3/etcdserver/api/membership"
	"go.etcd.io/etcd/server/v3/etcdserver/api/v2error"
	"go.etcd.io/etcd/server/v3/etcdserver/api/v2http/httptypes"
	"log"
	"net/http"
	"strconv"
	"strings"
)

func (s *_Server) peerMembersHandler(w http.ResponseWriter, r *http.Request) {
	//s.mu.RLock()
	//defer s.mu.RUnlock()
	if !allowMethod(w, r, "GET") {
		return
	}
	cluster := s.server.Cluster()

	w.Header().Set("X-Etcd-Cluster-ID", cluster.ID().String())

	if r.URL.Path != peerMembersPath {
		http.Error(w, "bad path", http.StatusBadRequest)
		return
	}
	ms := s.server.Cluster().Members()
	w.Header().Set("Content-Type", "application/json")
	if err := json.NewEncoder(w).Encode(ms); err != nil {
		log.Printf("failed to encode membership members: %s", err.Error())
	}
}

func (s *_Server) peerMemberPromoteHandler(w http.ResponseWriter, r *http.Request) {
	//s.mu.RLock()
	//defer s.mu.RUnlock()

	if !allowMethod(w, r, "POST") {
		return
	}
	w.Header().Set("X-Etcd-Cluster-ID", s.server.Cluster().ID().String())

	if !strings.HasPrefix(r.URL.Path, peerMemberPromotePrefix) {
		http.Error(w, "bad path", http.StatusBadRequest)
		return
	}
	idStr := strings.TrimPrefix(r.URL.Path, peerMemberPromotePrefix)
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		http.Error(w, fmt.Sprintf("member %s not found in cluster", idStr), http.StatusNotFound)
		return
	}

	resp, err := s.server.PromoteMember(r.Context(), id)
	if err != nil {
		switch err {
		case membership.ErrIDNotFound:
			http.Error(w, err.Error(), http.StatusNotFound)
		case membership.ErrMemberNotLearner:
			http.Error(w, err.Error(), http.StatusPreconditionFailed)
		case etcdserver.ErrLearnerNotReady:
			http.Error(w, err.Error(), http.StatusPreconditionFailed)
		default:
			WriteError(w, r, err)
		}
		log.Printf("failed to promote a member: id(%s)", types.ID(id).String())

		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	if err := json.NewEncoder(w).Encode(resp); err != nil {
		log.Printf("failed to encode members response: %s", err.Error())
	}
}
func WriteError(w http.ResponseWriter, r *http.Request, err error) {
	if err == nil {
		return
	}
	var v interface{}
	switch e := v.(type) {
	case *v2error.Error:
		_ = e.WriteTo(w)

	case *httptypes.HTTPError:
		if et := e.WriteTo(w); et != nil {
			log.Printf("failed to write v2 HTTP error, remote-addr: %s, internal-server-error: %s, %s", r.RemoteAddr, err.Error(), et.Error())
		}

	default:
		switch {
		case errors.Is(err, etcdserver.ErrTimeoutDueToLeaderFail), errors.Is(err, etcdserver.ErrTimeoutDueToConnectionLost), errors.Is(err, etcdserver.ErrNotEnoughStartedMembers), errors.Is(err, etcdserver.ErrUnhealthy):
			log.Printf("v2 response error, remote-addr: %s, internal-server-error: %s", r.RemoteAddr, err.Error())
		default:
			log.Printf("unexpected v2 response error, remote-addr: %s, internal-server-error: %s", r.RemoteAddr, err.Error())
		}

		herr := httptypes.NewHTTPError(http.StatusInternalServerError, "Internal Server Error")
		if et := herr.WriteTo(w); et != nil {
			log.Printf("failed to write v2 HTTP error, remote-addr: %s, internal-server-error: %s, %s", r.RemoteAddr, err.Error(), et.Error())

		}
	}
}
func allowMethod(w http.ResponseWriter, r *http.Request, m string) bool {
	if m == r.Method {
		return true
	}
	w.Header().Set("Allow", m)
	http.Error(w, "Method Not Allowed", http.StatusMethodNotAllowed)
	return false
}
