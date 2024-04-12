// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: master.proto

package service

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// MasterDispatcherClient is the client API for MasterDispatcher service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MasterDispatcherClient interface {
	Listen(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (MasterDispatcher_ListenClient, error)
}

type masterDispatcherClient struct {
	cc grpc.ClientConnInterface
}

func NewMasterDispatcherClient(cc grpc.ClientConnInterface) MasterDispatcherClient {
	return &masterDispatcherClient{cc}
}

func (c *masterDispatcherClient) Listen(ctx context.Context, in *EmptyRequest, opts ...grpc.CallOption) (MasterDispatcher_ListenClient, error) {
	stream, err := c.cc.NewStream(ctx, &MasterDispatcher_ServiceDesc.Streams[0], "/service.MasterDispatcher/Listen", opts...)
	if err != nil {
		return nil, err
	}
	x := &masterDispatcherListenClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type MasterDispatcher_ListenClient interface {
	Recv() (*Event, error)
	grpc.ClientStream
}

type masterDispatcherListenClient struct {
	grpc.ClientStream
}

func (x *masterDispatcherListenClient) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MasterDispatcherServer is the server API for MasterDispatcher service.
// All implementations must embed UnimplementedMasterDispatcherServer
// for forward compatibility
type MasterDispatcherServer interface {
	Listen(*EmptyRequest, MasterDispatcher_ListenServer) error
	mustEmbedUnimplementedMasterDispatcherServer()
}

// UnimplementedMasterDispatcherServer must be embedded to have forward compatible implementations.
type UnimplementedMasterDispatcherServer struct {
}

func (UnimplementedMasterDispatcherServer) Listen(*EmptyRequest, MasterDispatcher_ListenServer) error {
	return status.Errorf(codes.Unimplemented, "method Listen not implemented")
}
func (UnimplementedMasterDispatcherServer) mustEmbedUnimplementedMasterDispatcherServer() {}

// UnsafeMasterDispatcherServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MasterDispatcherServer will
// result in compilation errors.
type UnsafeMasterDispatcherServer interface {
	mustEmbedUnimplementedMasterDispatcherServer()
}

func RegisterMasterDispatcherServer(s grpc.ServiceRegistrar, srv MasterDispatcherServer) {
	s.RegisterService(&MasterDispatcher_ServiceDesc, srv)
}

func _MasterDispatcher_Listen_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(EmptyRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(MasterDispatcherServer).Listen(m, &masterDispatcherListenServer{stream})
}

type MasterDispatcher_ListenServer interface {
	Send(*Event) error
	grpc.ServerStream
}

type masterDispatcherListenServer struct {
	grpc.ServerStream
}

func (x *masterDispatcherListenServer) Send(m *Event) error {
	return x.ServerStream.SendMsg(m)
}

// MasterDispatcher_ServiceDesc is the grpc.ServiceDesc for MasterDispatcher service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MasterDispatcher_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.MasterDispatcher",
	HandlerType: (*MasterDispatcherServer)(nil),
	Methods:     []grpc.MethodDesc{},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "Listen",
			Handler:       _MasterDispatcher_Listen_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "master.proto",
}

// MasterEventsClient is the client API for MasterEvents service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MasterEventsClient interface {
	Send(ctx context.Context, in *Event, opts ...grpc.CallOption) (*EmptyResponse, error)
	SendStream(ctx context.Context, opts ...grpc.CallOption) (MasterEvents_SendStreamClient, error)
}

type masterEventsClient struct {
	cc grpc.ClientConnInterface
}

func NewMasterEventsClient(cc grpc.ClientConnInterface) MasterEventsClient {
	return &masterEventsClient{cc}
}

func (c *masterEventsClient) Send(ctx context.Context, in *Event, opts ...grpc.CallOption) (*EmptyResponse, error) {
	out := new(EmptyResponse)
	err := c.cc.Invoke(ctx, "/service.MasterEvents/Send", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *masterEventsClient) SendStream(ctx context.Context, opts ...grpc.CallOption) (MasterEvents_SendStreamClient, error) {
	stream, err := c.cc.NewStream(ctx, &MasterEvents_ServiceDesc.Streams[0], "/service.MasterEvents/SendStream", opts...)
	if err != nil {
		return nil, err
	}
	x := &masterEventsSendStreamClient{stream}
	return x, nil
}

type MasterEvents_SendStreamClient interface {
	Send(*Event) error
	CloseAndRecv() (*EmptyResponse, error)
	grpc.ClientStream
}

type masterEventsSendStreamClient struct {
	grpc.ClientStream
}

func (x *masterEventsSendStreamClient) Send(m *Event) error {
	return x.ClientStream.SendMsg(m)
}

func (x *masterEventsSendStreamClient) CloseAndRecv() (*EmptyResponse, error) {
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	m := new(EmptyResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MasterEventsServer is the server API for MasterEvents service.
// All implementations must embed UnimplementedMasterEventsServer
// for forward compatibility
type MasterEventsServer interface {
	Send(context.Context, *Event) (*EmptyResponse, error)
	SendStream(MasterEvents_SendStreamServer) error
	mustEmbedUnimplementedMasterEventsServer()
}

// UnimplementedMasterEventsServer must be embedded to have forward compatible implementations.
type UnimplementedMasterEventsServer struct {
}

func (UnimplementedMasterEventsServer) Send(context.Context, *Event) (*EmptyResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Send not implemented")
}
func (UnimplementedMasterEventsServer) SendStream(MasterEvents_SendStreamServer) error {
	return status.Errorf(codes.Unimplemented, "method SendStream not implemented")
}
func (UnimplementedMasterEventsServer) mustEmbedUnimplementedMasterEventsServer() {}

// UnsafeMasterEventsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MasterEventsServer will
// result in compilation errors.
type UnsafeMasterEventsServer interface {
	mustEmbedUnimplementedMasterEventsServer()
}

func RegisterMasterEventsServer(s grpc.ServiceRegistrar, srv MasterEventsServer) {
	s.RegisterService(&MasterEvents_ServiceDesc, srv)
}

func _MasterEvents_Send_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Event)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MasterEventsServer).Send(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.MasterEvents/Send",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MasterEventsServer).Send(ctx, req.(*Event))
	}
	return interceptor(ctx, in, info, handler)
}

func _MasterEvents_SendStream_Handler(srv interface{}, stream grpc.ServerStream) error {
	return srv.(MasterEventsServer).SendStream(&masterEventsSendStreamServer{stream})
}

type MasterEvents_SendStreamServer interface {
	SendAndClose(*EmptyResponse) error
	Recv() (*Event, error)
	grpc.ServerStream
}

type masterEventsSendStreamServer struct {
	grpc.ServerStream
}

func (x *masterEventsSendStreamServer) SendAndClose(m *EmptyResponse) error {
	return x.ServerStream.SendMsg(m)
}

func (x *masterEventsSendStreamServer) Recv() (*Event, error) {
	m := new(Event)
	if err := x.ServerStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// MasterEvents_ServiceDesc is the grpc.ServiceDesc for MasterEvents service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MasterEvents_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.MasterEvents",
	HandlerType: (*MasterEventsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Send",
			Handler:    _MasterEvents_Send_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SendStream",
			Handler:       _MasterEvents_SendStream_Handler,
			ClientStreams: true,
		},
	},
	Metadata: "master.proto",
}
