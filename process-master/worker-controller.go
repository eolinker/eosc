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
	Start([]int)
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
	restartChan       chan chan []int
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
		restartChan:       make(chan chan []int, 1),
	}
}

func (wc *WorkerController) Stop() {
	wc.locker.Lock()
	defer wc.locker.Unlock()

	if wc.isStop {
		return
	}
	close(wc.checkClose)
	close(wc.restartChan)
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
func (wc *WorkerController) Start(initPorts []int) {

	wc.NewWorker()

	go func() {
		t := time.NewTicker(time.Second / 4)
		in := &service.WorkerHelloRequest{
			Hello: "hello",
		}
		next := time.NewTimer(time.Second)
		next.Stop()
		var last = initPorts
		defer next.Stop()
		defer t.Stop()
		for {
			select {
			case <-t.C:
				client := wc.getClient()
				if client != nil {
					response, err := client.Ping(context.TODO(), in)
					if err != nil {
						log.Debug("ping worker controller error: ", err)
						continue
					}
					ports := sortAndSet(response.Resource.Port)

					if !equal(last, ports) {
						log.Debug("sort ports: ", ports, "last ports: ", last)
						last = ports
						next.Reset(time.Second)
					}
				}
			case <-next.C:
				{
					log.Debug("reset traffic:", last)
					isCreate, err := wc.trafficController.Reset(last)
					if err != nil {
						log.Debug("reset ports error: ", err, " last ports: ", last, " isCreate: ", isCreate)
						continue
					}
					if isCreate {
						wc.NewWorker()
					}
				}
			case <-wc.checkClose:
				return
			case cback, ok := <-wc.restartChan:
				if ok {
					cback <- last
					close(cback)
				}

				return
				//next.Reset(time.Second * 1)
			}
		}

	}()
}

func (wc *WorkerController) Restart() {
	//wc.trafficController.Reset(nil)
	cback := make(chan []int, 1)
	wc.restartChan <- cback

	last := <-cback
	wc.Start(last)

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

	for {
		select {
		case <-ticker.C:
			wc.current.createClient()
			_, err := wc.current.Ping(context.TODO(), &service.WorkerHelloRequest{Hello: "hello"})
			if err != nil {
				log.Debug("work controller ping: ", err)
				continue
			}

			return nil
		}
	}
	return nil
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

func equal(v1, v2 []int) bool {
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
func sortAndSet(vs []int32) []int {
	if len(vs) == 0 {
		return nil
	}

	m := make(map[int]int)
	for _, v := range vs {
		m[int(v)] = 1
	}
	rs := make([]int, 0, len(m))
	for v := range m {
		rs = append(rs, v)
	}
	sort.Ints(rs)
	return rs
}
