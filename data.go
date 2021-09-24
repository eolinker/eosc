package eosc

import (
	"os"
)

type IDataMarshaler interface {
	Encode(startIndex int) ([]byte, []*os.File, error)
}
