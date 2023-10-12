package process_worker

import (
	"context"
	"fmt"
	"io"
	"sync"
	"time"

	"github.com/eolinker/eosc/common/bean"
	"github.com/eolinker/eosc/professions"
	"github.com/eolinker/eosc/setting"
	"github.com/eolinker/eosc/variable"

	"github.com/eolinker/eosc"

	"github.com/eolinker/eosc/process-worker/workers"

	grpc_unixsocket "github.com/eolinker/eosc/grpc-unixsocket"
	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc/extends"

	"google.golang.org/grpc"

	"github.com/eolinker/eosc/log"
	"github.com/eolinker/eosc/service"
)

type WorkerServer struct {
	ctx               context.Context
	cancel            context.CancelFunc
	workers           workers.IWorkers
	professionManager professions.IProfessions
	settings          eosc.ISettings
	variableManager   eosc.IVariable
	masterPid         int
	onceInit          sync.Once
	initHandler       []func()
}

func NewWorkerServer(masterPid int, extends extends.IExtenderRegister, initHandlers ...func()) (*WorkerServer, error) {
	defer utils.TimeSpend("NewWorkerServer")()
	ctx, cancel := context.WithCancel(context.Background())
	ws := &WorkerServer{
		ctx:               ctx,
		cancel:            cancel,
		masterPid:         masterPid,
		professionManager: professions.NewProfessions(extends),
		initHandler:       initHandlers,
		variableManager:   variable.NewVariables(nil),
		settings:          setting.GetSettings(),
	}

	ws.workers = workers.NewWorkerManager(ws.professionManager)
	var iw eosc.IWorkers = ws.workers
	bean.Injection(&iw)
	ws.listenMaster()
	return ws, nil
}

func (ws *WorkerServer) Stop() {
	ws.cancel()
}

func (ws *WorkerServer) listenMaster() {
	conn, client, err := ws.createClient()
	if err == nil {

		go ws.listen(conn, client)
	} else {
		ws.retryConn()
	}
}

func (ws *WorkerServer) retryConn() {
	left, right := 1, 1
	ticker := time.NewTicker(time.Duration(left*5) * time.Second)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			{
				conn, c, err := ws.createClient()
				if err != nil {
					log.Error("create conn error: ", err)
					left, right = right, left+right
					ticker.Reset(time.Duration(left*5) * time.Second)
					continue
				}

				go ws.listen(conn, c)
				return
			}
		}
	}
}

func (ws *WorkerServer) createClient() (*grpc.ClientConn, service.MasterDispatcher_ListenClient, error) {
	addr := service.ServerAddr(ws.masterPid, eosc.ProcessMaster)
	conn, err := grpc_unixsocket.Connect(addr)
	if err != nil {
		return nil, nil, fmt.Errorf("connect master grpc addr error: %w,pid: %d\n", err, ws.masterPid)
	}

	client := service.NewMasterDispatcherClient(conn)
	opts := []grpc.CallOption{
		grpc.MaxCallRecvMsgSize(1024 * 1024 * 1024),
		grpc.MaxCallSendMsgSize(1024 * 1024 * 1024),
	}
	c, err := client.Listen(ws.ctx, &service.EmptyRequest{}, opts...)
	if err != nil {
		return nil, nil, fmt.Errorf("listen master service error: %w,pid: %d\n", err, ws.masterPid)
	}
	return conn, c, nil
}

func (ws *WorkerServer) listen(conn *grpc.ClientConn, c service.MasterDispatcher_ListenClient) {
	log.Debug("start listen")
	defer conn.Close()
	for {
		event, err := c.Recv()

		if err != nil {
			log.Error("recv:", err)
			if err == io.EOF {
				log.Debug("listen closed... ", err)
				return
			}
			ws.retryConn()
			return
		}
		log.Debug("recv:", event.String())
		switch event.Command {
		case eosc.EventInit, eosc.EventReset:
			{
				err := ws.resetEvent(event.Data)
				if err != nil {
					log.Error("reset server error: ", err)
					continue
				}
			}
		case eosc.EventSet:
			{
				ws.setEvent(event.Namespace, event.Key, event.Data)
			}
		case eosc.EventDel:
			{
				ws.delEvent(event.Namespace, event.Key)
			}
		}
	}
	log.Debug("stop listen")
}
