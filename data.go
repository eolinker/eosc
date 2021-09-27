package eosc

import (
	"os"
)

type IDataMarshaller interface {
	Encode(startIndex int) ([]byte, []*os.File, error)
}
