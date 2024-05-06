package model

import "github.com/eolinker/eosc/process-admin/cmd/proto"

type Object []byte

func (o *Object) UnmarshalBinary(data []byte) error {
	*o = data
	return nil
}

func (o *Object) String() string {
	return string(*o)
}

func (o *Object) Scan(v any) error {
	return proto.Scan(*o, v)
}
