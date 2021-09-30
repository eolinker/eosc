package raft

// 业务处理，根据实际需求更改service，service是外层的业务对象

type IService interface {
	// CommitHandler 节点commit信息前的处理
	CommitHandler(data []byte) (err error)

	// ProcessHandler 节点propose信息前的处理,leader发起或者未启用集群时使用
	//ProcessHandler(body interface{}) (data []byte, err error)
	// ProcessHandler 转发到leader时的处理
	ProcessDataHandler(body []byte) (object interface{}, data []byte, err error)

	// GetInit 集群初始化时的将service缓存中的信息进行打包处理,只会在切换集群模式的时候调用一次
	GetInit() (data []byte, err error)

	// ResetSnap 读取快照，用于恢复service数据
	ResetSnap(data []byte) (err error)

	// GetSnapshot 生成快照，用于快照文件的生成
	GetSnapshot() (data []byte, err error)

	SetRaft(raft IRaftSender)
}

type IRaftSender interface {
	Send(msg []byte) (interface{}, error)
}

type IRaftService interface {
	IService
}

//
//type raftDemo struct {
//	ser       IService
//	isCluster bool
//	isLeader  bool
//}
//type Msg struct {
//	Namespace string                 `json:"namespace"`
//	Cmd       string                 `json:"cmd"`
//	Body      map[string]interface{} `json:"body"`
//}
//
//func (r *raftDemo) Encode(namespace, cmd string, body map[string]interface{}) ([]byte, error) {
//	msg := &Msg{
//		Namespace: namespace,
//		Cmd:       cmd,
//		Body:      body,
//	}
//
//	return json.Marshal(msg)
//}
//
//func (r *raftDemo) Decode(data []byte) (string, string, map[string]interface{}, error) {
//	m := new(Msg)
//	err := json.Unmarshal(data, m)
//	if err != nil {
//		return "", "", nil, err
//	}
//
//	return m.Namespace, m.Cmd, m.Body, nil
//
//}
//
//func (r *raftDemo) Send(namespace, cmd string, body map[string]interface{}) error {
//	if r.isCluster {
//		if r.isLeader {
//			msg, err := r.ser.ProcessHandler(namespace, cmd, body)
//			if err != nil {
//				return err
//			}
//			return r.process(addr, "do", msg)
//		} else {
//			msgRaft, err := r.Encode(namespace, cmd, body)
//			if err != nil {
//				return err
//			}
//			return r.sendToLeader(addr, msgRaft)
//		}
//
//	} else {
//		msg, err := r.ser.ProcessHandler(namespace, cmd, body)
//		if err != nil {
//			return err
//		}
//
//		return r.ser.CommitHandler(namespace, msg)
//
//	}
//}
//
//func (r *raftDemo) HanderHttpSend(w http.ResponseWriter, req *http.Request) {
//	if r.isCluster {
//		w.Write("error not cluster")
//		return
//	}
//	if !r.isLeader {
//		w.Write("error not leader") // or send to leader
//		return
//	}
//
//	err := r.doHandler(req.Body)
//	req.Body.Close()
//	if err != nil {
//		writeError(w, "504", "error", err.Error())
//		return
//	}
//	eosc.WriteSuccessNoData(w, "")
//}
//
//func (r *raftDemo) doHandler(reader io.Reader) error {
//	msgRaft, err := ioutil.ReadAll(reader)
//	if err != nil {
//		return err
//	}
//	namespace, cmd, body, err := r.Decode(msgRaft)
//	if err != nil {
//		return err
//	}
//	msg, err := r.ser.ProcessHandler(namespace, cmd, body)
//	if err != nil {
//
//		return err
//	}
//
//	return r.process(addr, "do", msg)
//
//}
