package store

import (
	"fmt"
	"github.com/eolinker/eosc"
)

var(
	ErrorReadOnly = fmt.Errorf("yaml :%w",eosc.ErrorStoreReadOnly)
)
