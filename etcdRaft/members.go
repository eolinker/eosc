package etcdRaft

import (
	"fmt"
	"go.etcd.io/etcd/client/pkg/v3/types"
	"go.etcd.io/etcd/server/v3/etcdserver/api/membership"
	"strconv"
	"time"
)

func (e *EtcdServer) members() map[string][]string {
	e.mu.RLock()
	defer e.mu.RUnlock()
	list := make(map[string][]string)
	members := e.server.Cluster().Members()
	for _, member := range members {
		m := member.Clone()
		if len(m.Name) > 0 {
			list[m.Name] = m.PeerURLs
			continue
		}
		list[strconv.FormatUint(uint64(m.ID), 10)] = m.PeerURLs
	}
	return list
}

func (e *EtcdServer) addMember(urls []string) (map[string][]string, string, error) {
	purls, err := types.NewURLs(urls)
	if err != nil {
		return nil, "", err
	}
	now := time.Now()
	member := membership.NewMember("", purls, defaultClusterName, &now)
	member.Name = fmt.Sprintf("%s_%d", defaultName, member.ID)
	members, err := func() ([]*membership.Member, error) {
		ctx, cancel := e.requestContext()
		defer cancel()
		return e.server.AddMember(ctx, *member)
	}()
	res := make(map[string][]string)
	for _, m := range members {
		res[m.Name] = m.PeerURLs
	}
	return res, member.Name, nil
}

func (e *EtcdServer) updateSelfInfo(peers []string) error {
	member := e.server.Cluster().Member(e.server.ID())
	member.PeerURLs = peers
	return e.updateMember(*member)
}

func (e *EtcdServer) updateMember(member membership.Member) error {
	ctx, cancel := e.requestContext()
	defer cancel()
	_, err := e.server.UpdateMember(ctx, member)
	return err
}

func (e *EtcdServer) removeMember(id uint64) error {
	ctx, cancel := e.requestContext()
	defer cancel()
	_, err := e.server.RemoveMember(ctx, id)
	return err
}
