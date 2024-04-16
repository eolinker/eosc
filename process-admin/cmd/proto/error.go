package proto

import "errors"

var (
	ErrorCantParseArrayToString = errors.New("can't parse array to string")
	ErrorCantParseArrayToBool   = errors.New("can't parse array to string")
	ErrorCantParseArrayToInt    = errors.New("can't parse array to int")
	ErrorCantParseArrayToFloat  = errors.New("can't parse array to float")
)
