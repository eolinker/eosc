package log

import (
	"io"
	"os"
)

func NewStdTransport(formatter Formatter) *Transporter {
	t:= NewTransport(os.Stderr, InfoLevel)
	t.SetFormatter(formatter)
	return t
}

func NewTransport(out io.Writer, level Level) *Transporter {
	return &Transporter{
		output: out,
		level:  level,
	}
}
