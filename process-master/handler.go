package process_master

import (
	"encoding/json"
	"net/http"
	"strings"
)

func (m *Master) EtcdNodesHandler(w http.ResponseWriter, r *http.Request) {
	nodes := m.etcdServer.Nodes()
	for _, node := range nodes {
		for i, surl := range node.Server {
			node.Server[i] = strings.TrimLeft(surl, "eosc://")
		}
	}
	json.NewEncoder(w).Encode(nodes)
}
func (m *Master) EtcdInfoHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(m.etcdServer.Status())
}
