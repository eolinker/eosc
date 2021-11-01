package extenders

import "errors"

var (
	ErrorNotExist       = errors.New("not exist")
	ErrorInvalidVersion = errors.New("invalid version")
	ErrorInvalidCommand = errors.New("invalid command")
)
