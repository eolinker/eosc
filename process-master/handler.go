package process_master

import (
	"encoding/json"
	"net/http"
)

func (m *Master) EtcdNodesHandler(w http.ResponseWriter, r *http.Request) {
	nodes := m.etcdServer.Nodes()

	json.NewEncoder(w).Encode(nodes)
}
func (m *Master) EtcdInfoHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(m.etcdServer.Status())
}
