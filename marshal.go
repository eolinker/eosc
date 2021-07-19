package eosc

type Marshaler interface {
	Marshal()(string,error)
}