package process_master

import (
	"encoding/json"
	"net/http"
)

func (m *Master) EtcdNodesHandler(w http.ResponseWriter, r *http.Request) {

}
func (m *Master) EtcdInfoHandler(w http.ResponseWriter, r *http.Request) {

	json.NewEncoder(w).Encode(m.etcdServer.Info())
}
func (m *Master) StatusHandler(w http.ResponseWriter, r *http.Request) {
	m.etcdServer.Info()
}
