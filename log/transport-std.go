package log

import (
	"io"
)

func NewTransport(out io.Writer, level Level) *Transporter {
	return &Transporter{
		output: out,
		level:  level,
	}
}
