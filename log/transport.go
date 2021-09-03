package log

import (
	"io"
	"sync/atomic"
)

type Transporter struct {
	output    io.Writer
	level     Level
	formatter Formatter
}

func (t *Transporter) Output() io.Writer {
	return t.output
}

func (t *Transporter) SetOutput(output io.Writer) {
	t.output = output
}

func (t *Transporter) SetFormatter(formatter Formatter) {
	t.formatter = formatter
}

//func (t *Transporter) SetClose(close func() error) {
//	t.close = close
//}

func (t *Transporter) Transport(entry *Entry) error {
	output := t.output
	if output == nil {
		return nil
	}
	if t.Level() >= entry.Level {
		data, err := t.formatter.Format(entry)
		if err != nil {
			return err
		}
		_, err = output.Write(data)
		return err
	}
	return nil
}

func (t *Transporter) Level() Level {
	return Level(atomic.LoadUint32((*uint32)(&t.level)))
}

// SetLevel sets the logger level.
func (t *Transporter) SetLevel(level Level) {
	atomic.StoreUint32((*uint32)(&t.level), uint32(level))
}
func (t *Transporter) Close() error {
	t.output = nil
	return nil
}
