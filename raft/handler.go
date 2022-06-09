package raft

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"
	"strconv"
	"strings"

	"github.com/eolinker/eosc"

	"github.com/go-basic/uuid"

	"github.com/eolinker/eosc/log"
)

func (rc *Node) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
	if rc.transportHandler != nil {
		rc.transportHandler.ServeHTTP(writer, request)
		return
	}
	http.NotFound(writer, request)

}

// genHandler http请求处理
func (rc *Node) genHandler() http.Handler {
	sm := http.NewServeMux()
	// 其他节点加入集群的处理
	sm.HandleFunc("/raft/node/join", rc.joinHandler)

	sm.HandleFunc("/raft/node/join/try", rc.joinTry)

	sm.HandleFunc("/raft/node/join/callback", rc.joinCallback)

	// 其他节点转发到leader的处理
	//sm.HandleFunc("/raft/node/propose", rc.proposeHandler)

	sm.Handle("/", rc.transport.Handler())
	return sm
}

//joinCallback 分配节点信息
func (rc *Node) joinCallback(w http.ResponseWriter, r *http.Request) {
	writeSuccessResult(w, "", &SNResponse{
		SN: rc.lastSN,
	})
	return
}

//assignNodeInfo 分配节点信息
func (rc *Node) joinTryPost(w http.ResponseWriter, r *http.Request) {
	if !rc.join {
		msg := "invalid operation"
		writeError(w, "110005", msg, msg)
		return
	}
	ip := r.URL.Query().Get("broadcast_ip")
	port := r.URL.Query().Get("broadcast_port")
	protocol := r.URL.Query().Get("protocol")
	address := r.URL.Query().Get("address")
	p, _ := strconv.Atoi(port)
	addr := fmt.Sprintf("%s://%s", protocol, ip)
	if p > 0 {
		addr = fmt.Sprintf("%s:%d", addr, p)
	}
	callbackSN, err := callbackSNRequest(addr)
	if err != nil {
		writeError(w, "110006", "fail to callback sn", err.Error())
		return
	}
	sn := buildSN(ip, port, protocol, address, strconv.Itoa(os.Getpid()), eosc.GetRealIP(r), rc.nodeKey)
	if callbackSN != sn {
		msg := fmt.Sprintf("invalid sn,callback sn is %s,record sn is %s", callbackSN, sn)
		writeError(w, "110007", msg, msg)
		return
	}
	writeSuccessResult(w, "", &JoinResponse{
		NodeSecret: &NodeSecret{
			ID:  rc.peers.Index() + 1,
			Key: uuid.New(),
		},
		Peer: rc.peers.GetAllPeers(),
	})
	return
}

//buildSN 生成sn码
func buildSN(args ...string) string {
	builder := strings.Builder{}
	for _, a := range args {
		builder.WriteString(a)
	}
	h := md5.New()
	h.Write([]byte(builder.String()))
	return hex.EncodeToString(h.Sum(nil)) // 输出加密结果
}

//joinTryGet 获取校验码
func (rc *Node) joinTryGet(w http.ResponseWriter, r *http.Request) {
	ip := r.URL.Query().Get("broadcast_ip")
	port := r.URL.Query().Get("broadcast_port")
	protocol := r.URL.Query().Get("protocol")
	address := r.URL.Query().Get("address")

	if !rc.join {
		err := rc.UpdateHostInfo(address)
		if err != nil {
			writeError(w, "110002", "fail to update host Info", err.Error())
			return
		}
	}

	writeSuccessResult(w, "", &SNResponse{
		SN: buildSN(ip, port, protocol, address, strconv.Itoa(os.Getpid()), eosc.GetRealIP(r), rc.nodeKey),
	})
	return
}

func (rc *Node) joinTry(w http.ResponseWriter, r *http.Request) {
	switch r.Method {
	case http.MethodGet:
		rc.joinTryGet(w, r)
	case http.MethodPost:
		rc.joinTryPost(w, r)
	default:
		w.Write([]byte("invalid method type"))
		w.WriteHeader(http.StatusUnsupportedMediaType)
		return
	}
}

// joinHandler 收到其他节点加入集群的处理
// 1、如果已经是集群模式，直接返回相关id，peer等信息方便处理
// 2、如果不是集群模式，先切换集群rc.changeCluster,再返回相关信息
// 3、该处理也可应用于集群节点crash后的重启
func (rc *Node) joinHandler(w http.ResponseWriter, r *http.Request) {
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		return
	}
	defer r.Body.Close()
	joinData, err := decodeJoinRequest(body)
	if err != nil {
		writeError(w, "110001", "fail to parse join data", err.Error())
		return
	}

	addr := fmt.Sprintf("%s://%s", joinData.Protocol, joinData.BroadcastIP)
	if joinData.BroadcastPort > 0 {
		addr = fmt.Sprintf("%s:%d", addr, joinData.BroadcastPort)
	}
	log.Infof("address %s apply join the cluster", addr)
	if id, exist := rc.peers.CheckExist(addr); exist {
		// 当前地址已经存在
		if id != joinData.NodeID {
			// ID错误
			writeError(w, "110004", "id and address do not match", "id and address do not match")
			return
		}
		writeSuccessResult(w, "", nil)
		return
	}

	node := &NodeInfo{
		NodeSecret: &NodeSecret{
			ID:  joinData.NodeID,
			Key: joinData.NodeKey,
		},
		BroadcastIP:   joinData.BroadcastIP,
		BroadcastPort: joinData.BroadcastPort,
		Addr:          addr,
		Protocol:      joinData.Protocol,
	}
	data, err := json.Marshal(node)
	if err != nil {
		panic(err)
	}
	// 发送新增节点请求
	err = rc.AddNode(node.ID, data)
	if err != nil {
		writeError(w, "110003", "fail to add config error", err.Error())
		return
	}
	writeSuccessResult(w, "", nil)
	return
}

func decodeSNResponse(data []byte) (*SNResponse, error) {
	snRequest := new(SNResponse)
	err := json.Unmarshal(data, snRequest)
	if err != nil {
		return nil, err
	}
	return snRequest, nil
}

func decodeJoinRequest(data []byte) (*JoinRequest, error) {
	joinData := new(JoinRequest)
	err := json.Unmarshal(data, joinData)
	if err != nil {
		return nil, err
	}
	return joinData, nil
}

func decodeJoinResponse(data []byte) (*JoinResponse, error) {
	res := new(JoinResponse)
	err := json.Unmarshal(data, res)
	if err != nil {
		return nil, err
	}
	return res, nil
}
