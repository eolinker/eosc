package filelog

import (
	"bufio"
	"bytes"
	"context"
	"errors"
	"github.com/eolinker/eosc/common/pool"
	"github.com/eolinker/eosc/log"
	"sync"
	"time"
)

// MaxBuffer buffer最大值
const MaxBuffer = 1024 * 500

var (
	bufferPool         = pool.New(func() *bytes.Buffer { return new(bytes.Buffer) })
	ErrorWriterNotOpen = errors.New("writer close")
)

// FileWriterByPeriod 文件周期写入
type FileWriterByPeriod struct {
	wC chan *bytes.Buffer

	enable     bool
	cancelFunc context.CancelFunc
	locker     sync.RWMutex
	wg         sync.WaitGroup
	resetChan  chan Config

	watcher *Watcher

	fileController *FileController
}

// NewFileWriteByPeriod 获取新的FileWriterByPeriod
func NewFileWriteByPeriod(cfg Config) *FileWriterByPeriod {
	w := &FileWriterByPeriod{
		locker:    sync.RWMutex{},
		wg:        sync.WaitGroup{},
		enable:    false,
		resetChan: make(chan Config),
	}

	w.Open(cfg)
	return w
}
func (w *FileWriterByPeriod) Watch() (*WatchHandler, error) {
	w.locker.RLock()
	defer w.locker.RUnlock()
	if !w.enable {
		return nil, ErrorWriterNotOpen
	}

	return w.watcher.Watch(), nil
}
func (w *FileWriterByPeriod) Reset(cfg Config) {

	w.resetChan <- cfg
}

// Open 打开
func (w *FileWriterByPeriod) Open(cfg Config) {

	w.locker.Lock()
	defer w.locker.Unlock()

	if w.enable {
		return
	}
	w.watcher = NewWatcher()
	ctx, cancel := context.WithCancel(context.Background())
	w.cancelFunc = cancel
	w.wC = make(chan *bytes.Buffer, 100)

	w.enable = true
	w.wg.Add(1)
	go func() {
		w.do(ctx, cfg)
		w.wg.Done()
	}()
}

// Close 关闭
func (w *FileWriterByPeriod) Close() {

	isClose := false
	w.locker.Lock()
	defer w.locker.Unlock()
	if !w.enable {
		return
	}

	if w.cancelFunc != nil {
		isClose = true
		w.cancelFunc()
		w.cancelFunc = nil
	}
	if w.watcher != nil {
		w.watcher.Close()
		w.watcher = nil
	}
	w.enable = false

	if isClose {
		w.wg.Wait()
	}
}

func (w *FileWriterByPeriod) Write(p []byte) (n int, err error) {

	l := len(p)

	if l == 0 {
		return
	}
	buffer := bufferPool.Get()
	buffer.Reset()
	buffer.Write(p)
	if p[l-1] != '\n' {
		buffer.WriteByte('\n')
	}
	w.locker.RLock()
	defer w.locker.RUnlock()
	if !w.enable {
		bufferPool.PUT(buffer)
		return l, nil
	}

	w.wC <- buffer
	w.watcher.write(p)
	return l, nil
}

func (w *FileWriterByPeriod) do(ctx context.Context, config Config) {
	lastConfig := config
	w.fileController = NewFileController(config)
	w.fileController.initFile()
	currFile, lastTag, e := w.fileController.openFile()
	if e != nil {
		log.Errorf("open log file:%s\n", e.Error())
		return
	}

	buf := bufio.NewWriter(currFile)
	t := time.NewTicker(time.Second * 5)
	defer t.Stop()
	tFlush := time.NewTimer(time.Second)

	resetFunc := func() {
		if lastTag != w.fileController.timeTag(time.Now()) {
			if buf.Buffered() > 0 {
				buf.Flush()
				tFlush.Reset(time.Second)
			}
			stat, err := currFile.Stat()
			if err != nil {
				return
			}
			if stat.Size() > 0 { // 如果当前文件为空,则忽略重名文件
				currFile.Close()

				w.fileController.history(lastTag)
				fnew, tag, err := w.fileController.openFile()

				if err != nil {
					return
				}
				lastTag = tag
				currFile = fnew
				buf.Reset(currFile)

				go w.fileController.dropHistory()
			}

		}
	}

	for {
		select {
		case <-ctx.Done():
			{
				for len(w.wC) > 0 {
					p := <-w.wC
					buf.Write(p.Bytes())
					bufferPool.PUT(p)
				}
				buf.Flush()
				currFile.Close()
				t.Stop()
				//w.wg.Done()
				return
			}

		case <-t.C:
			{

				resetFunc()

			}
		case <-tFlush.C:
			{
				if buf.Buffered() > 0 {
					_ = buf.Flush()
				}
				tFlush.Reset(time.Second)
			}
		case p := <-w.wC:
			{
				_, _ = buf.Write(p.Bytes())
				bufferPool.PUT(p)

				if buf.Buffered() > MaxBuffer {
					_ = buf.Flush()
				}
				tFlush.Reset(time.Second)
			}
		case cfg, ok := <-w.resetChan:
			{
				if ok && lastConfig.IsUpdate(&cfg) {
					lastConfig = cfg
					w.fileController = NewFileController(cfg)
					resetFunc()
				}
			}
		}
	}
}
