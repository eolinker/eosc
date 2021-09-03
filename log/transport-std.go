package log

import (
	"io"
	"os"
)

func NewStdTransport(formatter Formatter) *Transporter {
	return NewTransport(os.Stdout, InfoLevel)
}

func NewTransport(out io.Writer, level Level) *Transporter {
	return &Transporter{
		output: out,
		level:  level,
	}
}
