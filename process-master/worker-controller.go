package process_master

import (
	"bytes"
	"context"
	"os"
	"sort"
	"sync"
	"time"

	"github.com/eolinker/eosc/utils"

	"github.com/eolinker/eosc/traffic"

	"github.com/eolinker/eosc/log"

	"github.com/eolinker/eosc"
	"github.com/eolinker/eosc/service"
)

var _ service.WorkerServiceClient = (*WorkerController)(nil)

type WorkerProcessController interface {
	Stop()
	NewWorker() error
	Start()
	Restart()
}
type WorkerController struct {
	locker            sync.Mutex
	dms               []eosc.IDataMarshaller
	current           *WorkerProcess
	expireWorkers     []*WorkerProcess
	trafficController traffic.IController
	isStop            bool
	checkClose        chan int
	portsChan         chan []int32
	localPorts        []int32
}

func NewWorkerController(trafficController traffic.IController, dms ...eosc.IDataMarshaller) *WorkerController {
	dmsAll := make([]eosc.IDataMarshaller, 0, len(dms)+1)
	dmsAll = append(dmsAll, trafficController)
	for _, v := range dms {
		dmsAll = append(dmsAll, v)
	}

	return &WorkerController{
		trafficController: trafficController,
		dms:               dmsAll,
		checkClose:        make(chan int, 1),
		portsChan:         make(chan []int32, 1),
	}
}

func (wc *WorkerController) Stop() {
	wc.locker.Lock()
	defer wc.locker.Unlock()

	if wc.isStop {
		return
	}
	close(wc.checkClose)
	close(wc.portsChan)
	wc.isStop = true
	if wc.current != nil {
		wc.current.Close()
		wc.expireWorkers = append(wc.expireWorkers, wc.current)
		wc.current = nil
	}
}
func (wc *WorkerController) check(w *WorkerProcess) {
	err := w.cmd.Wait()
	if err != nil {
		log.Warn("worker exit:", err)
	}

	if wc.getClient() == w {
		err := wc.NewWorker()
		if err != nil {
			log.Error("worker create:", err)
		}
	} else {

		for i, v := range wc.expireWorkers {
			if v == w {
				wc.expireWorkers = append(wc.expireWorkers[:i], wc.expireWorkers[i+1:]...)
			}
		}
	}
}
func (wc *WorkerController) Start() {

	wc.NewWorker()

	go func() {

		next := time.NewTimer(time.Second)
		next.Stop()
		defer next.Stop()
		for {
			select {
			case <-next.C:
				{
					response, err := wc.current.Ping(context.TODO(), &service.WorkerHelloRequest{Hello: "hello"})
					if err != nil {
						continue
					}
					ps, psInt := sortAndSet(response.Resource.Port)
					if equal(ps, wc.localPorts) {
						continue
					}
					log.Debug("reset traffic:", ps)
					isCreate, err := wc.trafficController.Reset(psInt)
					if err != nil {
						log.Debug("reset ports error: ", err, " last ports: ", ps, " isCreate: ", isCreate)
						continue
					}
					wc.localPorts = ps
					if isCreate {
						wc.NewWorker()
					} else {
						in := &service.WorkerRefreshRequest{
							Ports: wc.localPorts,
						}
						_, err := wc.Refresh(context.TODO(), in)
						if err != nil {
							log.Debug("ping worker controller error: ", err)
							continue
						}
					}
				}
			case <-wc.checkClose:
				return
			case _, ok := <-wc.portsChan:
				if ok {
					//last = ports
					next.Reset(time.Second)
				}
			}
		}

	}()
}

func (wc *WorkerController) Restart() {
	wc.NewWorker()

}
func (wc *WorkerController) NewWorker() error {

	wc.locker.Lock()
	defer wc.locker.Unlock()
	err := wc.new()
	if err != nil {
		log.Warn("new worker:", err)
		return err
	}

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()
	defer utils.Timeout("wait worker process start:")()
	wc.current.createClient()
	for {
		_, err := wc.current.Ping(context.TODO(), &service.WorkerHelloRequest{Hello: "hello"})
		if err == nil {
			log.Debug("work controller ping: ", err)
			wc.portsChan <- []int32{}
			return nil
		}
		<-ticker.C
	}

}

func (wc *WorkerController) new() error {
	log.Debug("create worker process start")
	buf := bytes.NewBuffer(nil)
	var fileAll []*os.File
	index := 3
	for _, dm := range wc.dms {
		data, files, err := dm.Encode(index)
		log.Debugf("encode:data[%d] file[%d]", len(data), len(files))
		if err != nil {
			log.Warn("create worker process fail:", err)
			return err
		}
		index += len(files)
		fileAll = append(fileAll, files...)
		buf.Write(data)
	}

	workerProcess, err := newWorkerProcess(buf, fileAll)
	if err != nil {
		log.Warn("new worker process:", err)
		return err
	}

	old := wc.current
	wc.current = workerProcess
	go wc.check(wc.current)

	if old != nil {
		old.Close()
	}

	return nil
}

func (wc *WorkerController) getClient() *WorkerProcess {
	wc.locker.Lock()
	defer wc.locker.Unlock()
	if wc.current == nil {
		return nil
	}
	wc.current.createClient()
	return wc.current
}

func equal(v1, v2 []int32) bool {
	if len(v1) != len(v2) {
		return false
	}

	for i, v := range v1 {
		if v != v2[i] {
			return false
		}
	}
	return true
}

type int32Slice []int32

func (x int32Slice) Len() int {
	return len(x)
}

func (x int32Slice) Less(i, j int) bool {
	return x[i] < x[j]
}

func (x int32Slice) Swap(i, j int) {
	x[i], x[j] = x[j], x[i]
}

func sortAndSet(vs []int32) ([]int32, []int) {
	if len(vs) == 0 {
		return nil, nil
	}

	m := make(map[int32]int32)
	for _, v := range vs {
		m[v] = 1
	}
	rs := make([]int32, 0, len(m))
	rsInt := make([]int, 0, len(m))
	for v := range m {
		rs = append(rs, v)
		rsInt = append(rsInt, int(v))
	}
	sort.Sort(int32Slice(rs))
	sort.Ints(rsInt)
	return rs, rsInt
}
