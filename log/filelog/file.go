package filelog

import (
	"fmt"
	"github.com/eolinker/eosc/log"
	"os"
	"path/filepath"
	"strings"
	"time"
)

type FileController struct {
	expire time.Duration
	dir    string
	file   string
	period LogPeriod
}

func NewFileController(config Config) *FileController {
	return &FileController{
		dir:    config.Dir,
		file:   strings.TrimSuffix(config.File, ".log"),
		period: config.Period,
		expire: config.Expire,
	}
}

func (w *FileController) timeTag(t time.Time) string {

	tag := t.Format(w.period.FormatLayout())

	return filepath.Join(w.dir, fmt.Sprintf("%s-%s.log", w.file, tag))
}
func (w *FileController) fileName() string {
	return filepath.Join(w.dir, fmt.Sprintf("%s.log", w.file))
}
func (w *FileController) history(history string) {

	path := w.fileName()
	os.Rename(path, history)

}

func (w *FileController) dropHistory() {

	expireTime := time.Now().Add(-w.expire)
	pathPatten := filepath.Join(w.dir, fmt.Sprintf("%s-*", w.file))
	files, err := filepath.Glob(pathPatten)
	if err == nil {
		for _, f := range files {
			if info, e := os.Stat(f); e == nil {

				if expireTime.After(info.ModTime()) {
					_ = os.Remove(f)
				}
			}

		}
	}
}

func (w *FileController) initFile() {
	err := os.MkdirAll(w.dir, 0666)
	if err != nil {
		log.Error(err)
	}
	path := w.fileName()
	nowHistoryName := w.timeTag(time.Now())
	if info, e := os.Stat(path); e == nil {

		timeTag := w.timeTag(info.ModTime())
		if timeTag != nowHistoryName {
			w.history(timeTag)
		}
	}

	w.dropHistory()

}

func (w *FileController) openFile() (*os.File, string, error) {
	path := w.fileName()
	nowTag := w.timeTag(time.Now())
	f, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
		return nil, "", err
	}
	return f, nowTag, err

}
