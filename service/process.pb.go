package service

import (
	config "github.com/eolinker/eosc/config"
	traffic "github.com/eolinker/eosc/traffic"
)

type ProcessLoadArg struct {
	Traffic    []*traffic.PbTraffic `protobuf:"bytes,1,rep,name=traffic,proto3" json:"traffic,omitempty"`
	ListensMsg config.ListenUrl     `protobuf:"bytes,2,opt,name=listensMsg,proto3" json:"listensMsg,omitempty"`
	Extends    map[string]string    `protobuf:"bytes,3,rep,name=extends,proto3" json:"extends,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}
