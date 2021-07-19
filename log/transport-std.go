package log

import (
	"io"
	"os"
)

func NewStdTransport(formatter Formatter) *Transporter {
	return NewTransport(os.Stdout, InfoLevel, formatter)
}

func NewTransport(out io.Writer, level Level, formatter Formatter) *Transporter {
	return &Transporter{
		output:    out,
		level:     level,
		formatter: formatter,
	}
}
