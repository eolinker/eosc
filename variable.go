package eosc

import (
	"errors"
	"reflect"
)

var (
	ErrorUnsupportedKind = errors.New("unsupported kind")
)

type IVariable interface {
	SetByNamespace(namespace string, variables map[string]string) error
	GetByNamespace(namespace string) (map[string]string, bool)
	SetVariablesById(id string, variables []string)
	RemoveRequire(id string)
	Unmarshal(buf []byte, typ reflect.Type) (interface{}, []string, error)
	Check(namespace string, variables map[string]string) ([]string, IVariable, error)
	Get(id string) (string, bool)
	Len() int
}
