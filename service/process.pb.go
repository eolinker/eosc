// Code generated by protoc-gen-go. DO NOT EDIT.
// versions:
// 	protoc-gen-go v1.28.0
// 	protoc        v3.19.4
// source: process.proto

package service

import (
	config "github.com/eolinker/eosc/config"
	traffic "github.com/eolinker/eosc/traffic"
	protoreflect "google.golang.org/protobuf/reflect/protoreflect"
	protoimpl "google.golang.org/protobuf/runtime/protoimpl"
	reflect "reflect"
	sync "sync"
)

const (
	// Verify that this generated code is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(20 - protoimpl.MinVersion)
	// Verify that runtime/protoimpl is sufficiently up-to-date.
	_ = protoimpl.EnforceVersion(protoimpl.MaxVersion - 20)
)

type ProcessLoadArg struct {
	state         protoimpl.MessageState
	sizeCache     protoimpl.SizeCache
	unknownFields protoimpl.UnknownFields

	Traffic    []*traffic.PbTraffic `protobuf:"bytes,1,rep,name=traffic,proto3" json:"traffic,omitempty"`
	ListensMsg *config.ListensMsg   `protobuf:"bytes,2,opt,name=listensMsg,proto3" json:"listensMsg,omitempty"`
	Extends    map[string]string    `protobuf:"bytes,3,rep,name=extends,proto3" json:"extends,omitempty" protobuf_key:"bytes,1,opt,name=key,proto3" protobuf_val:"bytes,2,opt,name=value,proto3"`
}

func (x *ProcessLoadArg) Reset() {
	*x = ProcessLoadArg{}
	if protoimpl.UnsafeEnabled {
		mi := &file_process_proto_msgTypes[0]
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		ms.StoreMessageInfo(mi)
	}
}

func (x *ProcessLoadArg) String() string {
	return protoimpl.X.MessageStringOf(x)
}

func (*ProcessLoadArg) ProtoMessage() {}

func (x *ProcessLoadArg) ProtoReflect() protoreflect.Message {
	mi := &file_process_proto_msgTypes[0]
	if protoimpl.UnsafeEnabled && x != nil {
		ms := protoimpl.X.MessageStateOf(protoimpl.Pointer(x))
		if ms.LoadMessageInfo() == nil {
			ms.StoreMessageInfo(mi)
		}
		return ms
	}
	return mi.MessageOf(x)
}

// Deprecated: Use ProcessLoadArg.ProtoReflect.Descriptor instead.
func (*ProcessLoadArg) Descriptor() ([]byte, []int) {
	return file_process_proto_rawDescGZIP(), []int{0}
}

func (x *ProcessLoadArg) GetTraffic() []*traffic.PbTraffic {
	if x != nil {
		return x.Traffic
	}
	return nil
}

func (x *ProcessLoadArg) GetListensMsg() *config.ListensMsg {
	if x != nil {
		return x.ListensMsg
	}
	return nil
}

func (x *ProcessLoadArg) GetExtends() map[string]string {
	if x != nil {
		return x.Extends
	}
	return nil
}

var File_process_proto protoreflect.FileDescriptor

var file_process_proto_rawDesc = []byte{
	0x0a, 0x0d, 0x70, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x12,
	0x07, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x1a, 0x0d, 0x74, 0x72, 0x61, 0x66, 0x66, 0x69,
	0x63, 0x2e, 0x70, 0x72, 0x6f, 0x74, 0x6f, 0x1a, 0x0c, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x2e,
	0x70, 0x72, 0x6f, 0x74, 0x6f, 0x22, 0xef, 0x01, 0x0a, 0x0e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73,
	0x73, 0x4c, 0x6f, 0x61, 0x64, 0x41, 0x72, 0x67, 0x12, 0x2c, 0x0a, 0x07, 0x74, 0x72, 0x61, 0x66,
	0x66, 0x69, 0x63, 0x18, 0x01, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x12, 0x2e, 0x73, 0x65, 0x72, 0x76,
	0x69, 0x63, 0x65, 0x2e, 0x50, 0x62, 0x54, 0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x52, 0x07, 0x74,
	0x72, 0x61, 0x66, 0x66, 0x69, 0x63, 0x12, 0x33, 0x0a, 0x0a, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e,
	0x73, 0x4d, 0x73, 0x67, 0x18, 0x02, 0x20, 0x01, 0x28, 0x0b, 0x32, 0x13, 0x2e, 0x73, 0x65, 0x72,
	0x76, 0x69, 0x63, 0x65, 0x2e, 0x4c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x73, 0x4d, 0x73, 0x67, 0x52,
	0x0a, 0x6c, 0x69, 0x73, 0x74, 0x65, 0x6e, 0x73, 0x4d, 0x73, 0x67, 0x12, 0x3e, 0x0a, 0x07, 0x65,
	0x78, 0x74, 0x65, 0x6e, 0x64, 0x73, 0x18, 0x03, 0x20, 0x03, 0x28, 0x0b, 0x32, 0x24, 0x2e, 0x73,
	0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x2e, 0x50, 0x72, 0x6f, 0x63, 0x65, 0x73, 0x73, 0x4c, 0x6f,
	0x61, 0x64, 0x41, 0x72, 0x67, 0x2e, 0x45, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x73, 0x45, 0x6e, 0x74,
	0x72, 0x79, 0x52, 0x07, 0x65, 0x78, 0x74, 0x65, 0x6e, 0x64, 0x73, 0x1a, 0x3a, 0x0a, 0x0c, 0x45,
	0x78, 0x74, 0x65, 0x6e, 0x64, 0x73, 0x45, 0x6e, 0x74, 0x72, 0x79, 0x12, 0x10, 0x0a, 0x03, 0x6b,
	0x65, 0x79, 0x18, 0x01, 0x20, 0x01, 0x28, 0x09, 0x52, 0x03, 0x6b, 0x65, 0x79, 0x12, 0x14, 0x0a,
	0x05, 0x76, 0x61, 0x6c, 0x75, 0x65, 0x18, 0x02, 0x20, 0x01, 0x28, 0x09, 0x52, 0x05, 0x76, 0x61,
	0x6c, 0x75, 0x65, 0x3a, 0x02, 0x38, 0x01, 0x42, 0x22, 0x5a, 0x20, 0x67, 0x69, 0x74, 0x68, 0x75,
	0x62, 0x2e, 0x63, 0x6f, 0x6d, 0x2f, 0x65, 0x6f, 0x6c, 0x69, 0x6e, 0x6b, 0x65, 0x72, 0x2f, 0x65,
	0x6f, 0x73, 0x63, 0x2f, 0x73, 0x65, 0x72, 0x76, 0x69, 0x63, 0x65, 0x62, 0x06, 0x70, 0x72, 0x6f,
	0x74, 0x6f, 0x33,
}

var (
	file_process_proto_rawDescOnce sync.Once
	file_process_proto_rawDescData = file_process_proto_rawDesc
)

func file_process_proto_rawDescGZIP() []byte {
	file_process_proto_rawDescOnce.Do(func() {
		file_process_proto_rawDescData = protoimpl.X.CompressGZIP(file_process_proto_rawDescData)
	})
	return file_process_proto_rawDescData
}

var file_process_proto_msgTypes = make([]protoimpl.MessageInfo, 2)
var file_process_proto_goTypes = []interface{}{
	(*ProcessLoadArg)(nil),    // 0: service.ProcessLoadArg
	nil,                       // 1: service.ProcessLoadArg.ExtendsEntry
	(*traffic.PbTraffic)(nil), // 2: service.PbTraffic
	(*config.ListensMsg)(nil), // 3: service.ListensMsg
}
var file_process_proto_depIdxs = []int32{
	2, // 0: service.ProcessLoadArg.traffic:type_name -> service.PbTraffic
	3, // 1: service.ProcessLoadArg.listensMsg:type_name -> service.ListensMsg
	1, // 2: service.ProcessLoadArg.extends:type_name -> service.ProcessLoadArg.ExtendsEntry
	3, // [3:3] is the sub-list for method output_type
	3, // [3:3] is the sub-list for method input_type
	3, // [3:3] is the sub-list for extension type_name
	3, // [3:3] is the sub-list for extension extendee
	0, // [0:3] is the sub-list for field type_name
}

func init() { file_process_proto_init() }
func file_process_proto_init() {
	if File_process_proto != nil {
		return
	}
	if !protoimpl.UnsafeEnabled {
		file_process_proto_msgTypes[0].Exporter = func(v interface{}, i int) interface{} {
			switch v := v.(*ProcessLoadArg); i {
			case 0:
				return &v.state
			case 1:
				return &v.sizeCache
			case 2:
				return &v.unknownFields
			default:
				return nil
			}
		}
	}
	type x struct{}
	out := protoimpl.TypeBuilder{
		File: protoimpl.DescBuilder{
			GoPackagePath: reflect.TypeOf(x{}).PkgPath(),
			RawDescriptor: file_process_proto_rawDesc,
			NumEnums:      0,
			NumMessages:   2,
			NumExtensions: 0,
			NumServices:   0,
		},
		GoTypes:           file_process_proto_goTypes,
		DependencyIndexes: file_process_proto_depIdxs,
		MessageInfos:      file_process_proto_msgTypes,
	}.Build()
	File_process_proto = out.File
	file_process_proto_rawDesc = nil
	file_process_proto_goTypes = nil
	file_process_proto_depIdxs = nil
}
