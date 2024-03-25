package proto

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"strconv"
)

type Reader struct {
	rd   *bufio.Reader
	_buf []byte
}

func NewReader(rd io.Reader) *Reader {
	return &Reader{
		rd:   bufio.NewReader(rd),
		_buf: make([]byte, 64),
	}
}
func (r *Reader) Buffered() int {
	return r.rd.Buffered()
}
func (r *Reader) Peek(n int) ([]byte, error) {
	return r.rd.Peek(n)
}
func (r *Reader) Reset(rd io.Reader) {
	r.rd.Reset(rd)
}
func (r *Reader) ReadMessage() (IMessage, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	if isNilReply(line) {
		return nil, Nil
	}
	switch line[0] {
	case ErrorReply:
		return nil, Error(line[1:])
	case IntReply, StatusReply:
		return MessageBase(line), nil
	case StringReply:
		value, err := r._readTmpBytesReply(line)
		if err != nil {
			return nil, err
		}
		return MessageString(value), nil
	case ArrayReply:

		n, err := parseArrayLen(line)
		if err != nil {
			return nil, err
		}
		if n == 0 {
			return &MessageArray{
				items: nil,
				err:   Nil,
			}, Nil
		}
		children := make([]IMessage, 0, n)

		for i := int64(0); i < n; i++ {
			child, err := r.ReadMessage()
			if err != nil {
				return nil, err
			}
			children = append(children, child)
		}
		return &MessageArray{
			items: children,
			err:   nil,
		}, nil
	}
	return nil, fmt.Errorf("apinto: invalid reply: %q", line)
}
func (r *Reader) ReadLine() ([]byte, error) {
	line, err := r.readLine()
	if err != nil {
		return nil, err
	}
	if isNilReply(line) {
		return nil, Nil
	}
	return line, nil
}
func (r *Reader) buf(n int) []byte {
	if n <= cap(r._buf) {
		return r._buf[:n]
	}
	d := n - cap(r._buf)
	r._buf = append(r._buf, make([]byte, d)...)
	return r._buf
}
func (r *Reader) _readTmpBytesReply(line []byte) ([]byte, error) {
	if isNilReply(line) {
		return nil, Nil
	}

	replyLen, err := strconv.Atoi(string(line[1:]))
	if err != nil {
		return nil, err
	}

	buf := r.buf(replyLen + 2)
	_, err = io.ReadFull(r.rd, buf)
	if err != nil {
		return nil, err
	}

	return buf[:replyLen], nil
}

// readLine that returns an error if:
//   - there is a pending read error;
//   - or line does not end with \r\n.
func (r *Reader) readLine() ([]byte, error) {
	b, err := r.rd.ReadSlice('\n')
	if err != nil {
		if !errors.Is(err, bufio.ErrBufferFull) {
			return nil, err
		}

		full := make([]byte, len(b))
		copy(full, b)

		b, err = r.rd.ReadBytes('\n')
		if err != nil {
			return nil, err
		}

		full = append(full, b...) //nolint:makezero
		b = full
	}
	if len(b) <= 2 || b[len(b)-1] != '\n' || b[len(b)-2] != '\r' {
		return nil, fmt.Errorf("apinto: invalid reply: %q", b)
	}
	return b[:len(b)-2], nil
}
func isNilReply(b []byte) bool {
	return len(b) == 3 &&
		(b[0] == StringReply || b[0] == ArrayReply) &&
		b[1] == '-' && b[2] == '1'
}
func parseArrayLen(line []byte) (int64, error) {
	if isNilReply(line) {
		return 0, Nil
	}
	return strconv.ParseInt(string(line[1:]), 10, 64)
}
