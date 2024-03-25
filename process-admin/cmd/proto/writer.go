package proto

import (
	"encoding"
	"encoding/json"
	"io"
	"strconv"
	"strings"
	"time"
)

type writer interface {
	io.Writer
	io.ByteWriter
}
type Writer struct {
	writer

	lenBuf []byte
	numBuf []byte
}

func NewWriter(wr writer) *Writer {
	return &Writer{
		writer: wr,

		lenBuf: make([]byte, 64),
		numBuf: make([]byte, 64),
	}
}

func (w *Writer) WriteArgs(args ...any) error {
	if err := w.WriteByte(ArrayReply); err != nil {
		return err
	}
	if err := w.writeLen(len(args)); err != nil {
		return err
	}
	for _, arg := range args {
		if err := w.WriteArg(arg); err != nil {
			return err
		}
	}
	return nil
}
func (w *Writer) WriteArg(v any) error {

	switch v := v.(type) {
	case nil:
		return w.status(nil)
	case string:
		if strings.ContainsAny(v, "\r\n") {
			return w.bytes([]byte(v))
		}
		return w.status([]byte(v))
	case []byte:
		return w.bytes(v)
	case int:
		return w.int(int64(v))
	case int8:
		return w.int(int64(v))
	case int16:
		return w.int(int64(v))
	case int32:
		return w.int(int64(v))
	case int64:
		return w.int(v)
	case uint:
		return w.uint(uint64(v))
	case uint8:
		return w.uint(uint64(v))
	case uint16:
		return w.uint(uint64(v))
	case uint32:
		return w.uint(uint64(v))
	case uint64:
		return w.uint(v)
	case float32:
		return w.float(float64(v))
	case float64:
		return w.float(v)
	case bool:
		if v {
			return w.int(1)
		}
		return w.int(0)
	case time.Time:
		w.numBuf = v.AppendFormat(w.numBuf[:0], time.RFC3339Nano)
		return w.status(w.numBuf)
	case time.Duration:
		return w.int(v.Nanoseconds())
	case encoding.BinaryMarshaler:
		b, err := v.MarshalBinary()
		if err != nil {
			return err
		}
		return w.bytes(b)
	default:
		bytes, err := json.Marshal(v)
		if err != nil {
			return err
		}
		return w.bytes(bytes)
	}
}
func (w *Writer) bytes(b []byte) error {
	if err := w.WriteByte(StringReply); err != nil {
		return err
	}

	if err := w.writeLen(len(b)); err != nil {
		return err
	}

	if _, err := w.Write(b); err != nil {
		return err
	}

	return w.crlf()
}

func (w *Writer) status(v []byte) error {
	err := w.WriteByte(StatusReply)
	if err != nil {
		return err
	}
	if len(v) > 0 {
		_, err = w.Write(v)

		if err != nil {
			return err
		}
	}
	return w.crlf()
}
func (w *Writer) writeLen(n int) error {
	w.lenBuf = strconv.AppendUint(w.lenBuf[:0], uint64(n), 10)
	w.lenBuf = append(w.lenBuf, '\r', '\n')
	_, err := w.Write(w.lenBuf)
	return err
}
func (w *Writer) crlf() error {
	if err := w.WriteByte('\r'); err != nil {
		return err
	}
	return w.WriteByte('\n')
}

func (w *Writer) uint(n uint64) error {
	w.numBuf = strconv.AppendUint(w.numBuf[:0], n, 10)
	return w.status(w.numBuf)
}

func (w *Writer) int(n int64) error {
	w.numBuf = strconv.AppendInt(w.numBuf[:0], n, 10)
	return w.status(w.numBuf)
}

func (w *Writer) float(f float64) error {
	w.numBuf = strconv.AppendFloat(w.numBuf[:0], f, 'f', -1, 64)
	return w.status(w.numBuf)
}
