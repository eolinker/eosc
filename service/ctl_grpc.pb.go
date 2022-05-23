// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.4
// source: ctl.proto

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

// CtiServiceClient is the client API for CtiService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type CtiServiceClient interface {
	Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error)
	Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*LeaveResponse, error)
	List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error)
	Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error)
	ExtendsInstall(ctx context.Context, in *ExtendsRequest, opts ...grpc.CallOption) (*ExtendsResponse, error)
	ExtendsUpdate(ctx context.Context, in *ExtendsRequest, opts ...grpc.CallOption) (*ExtendsResponse, error)
	ExtendsUninstall(ctx context.Context, in *ExtendsRequest, opts ...grpc.CallOption) (*ExtendsUninstallResponse, error)
}

type ctiServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewCtiServiceClient(cc grpc.ClientConnInterface) CtiServiceClient {
	return &ctiServiceClient{cc}
}

func (c *ctiServiceClient) Join(ctx context.Context, in *JoinRequest, opts ...grpc.CallOption) (*JoinResponse, error) {
	out := new(JoinResponse)
	err := c.cc.Invoke(ctx, "/service.CtiService/Join", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ctiServiceClient) Leave(ctx context.Context, in *LeaveRequest, opts ...grpc.CallOption) (*LeaveResponse, error) {
	out := new(LeaveResponse)
	err := c.cc.Invoke(ctx, "/service.CtiService/Leave", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ctiServiceClient) List(ctx context.Context, in *ListRequest, opts ...grpc.CallOption) (*ListResponse, error) {
	out := new(ListResponse)
	err := c.cc.Invoke(ctx, "/service.CtiService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ctiServiceClient) Info(ctx context.Context, in *InfoRequest, opts ...grpc.CallOption) (*InfoResponse, error) {
	out := new(InfoResponse)
	err := c.cc.Invoke(ctx, "/service.CtiService/Info", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ctiServiceClient) ExtendsInstall(ctx context.Context, in *ExtendsRequest, opts ...grpc.CallOption) (*ExtendsResponse, error) {
	out := new(ExtendsResponse)
	err := c.cc.Invoke(ctx, "/service.CtiService/ExtendsInstall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ctiServiceClient) ExtendsUpdate(ctx context.Context, in *ExtendsRequest, opts ...grpc.CallOption) (*ExtendsResponse, error) {
	out := new(ExtendsResponse)
	err := c.cc.Invoke(ctx, "/service.CtiService/ExtendsUpdate", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *ctiServiceClient) ExtendsUninstall(ctx context.Context, in *ExtendsRequest, opts ...grpc.CallOption) (*ExtendsUninstallResponse, error) {
	out := new(ExtendsUninstallResponse)
	err := c.cc.Invoke(ctx, "/service.CtiService/ExtendsUninstall", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// CtiServiceServer is the server API for CtiService service.
// All implementations must embed UnimplementedCtiServiceServer
// for forward compatibility
type CtiServiceServer interface {
	Join(context.Context, *JoinRequest) (*JoinResponse, error)
	Leave(context.Context, *LeaveRequest) (*LeaveResponse, error)
	List(context.Context, *ListRequest) (*ListResponse, error)
	Info(context.Context, *InfoRequest) (*InfoResponse, error)
	ExtendsInstall(context.Context, *ExtendsRequest) (*ExtendsResponse, error)
	ExtendsUpdate(context.Context, *ExtendsRequest) (*ExtendsResponse, error)
	ExtendsUninstall(context.Context, *ExtendsRequest) (*ExtendsUninstallResponse, error)
	mustEmbedUnimplementedCtiServiceServer()
}

// UnimplementedCtiServiceServer must be embedded to have forward compatible implementations.
type UnimplementedCtiServiceServer struct {
}

func (UnimplementedCtiServiceServer) Join(context.Context, *JoinRequest) (*JoinResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Join not implemented")
}
func (UnimplementedCtiServiceServer) Leave(context.Context, *LeaveRequest) (*LeaveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Leave not implemented")
}
func (UnimplementedCtiServiceServer) List(context.Context, *ListRequest) (*ListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedCtiServiceServer) Info(context.Context, *InfoRequest) (*InfoResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Info not implemented")
}
func (UnimplementedCtiServiceServer) ExtendsInstall(context.Context, *ExtendsRequest) (*ExtendsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExtendsInstall not implemented")
}
func (UnimplementedCtiServiceServer) ExtendsUpdate(context.Context, *ExtendsRequest) (*ExtendsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExtendsUpdate not implemented")
}
func (UnimplementedCtiServiceServer) ExtendsUninstall(context.Context, *ExtendsRequest) (*ExtendsUninstallResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExtendsUninstall not implemented")
}
func (UnimplementedCtiServiceServer) mustEmbedUnimplementedCtiServiceServer() {}

// UnsafeCtiServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to CtiServiceServer will
// result in compilation errors.
type UnsafeCtiServiceServer interface {
	mustEmbedUnimplementedCtiServiceServer()
}

func RegisterCtiServiceServer(s grpc.ServiceRegistrar, srv CtiServiceServer) {
	s.RegisterService(&CtiService_ServiceDesc, srv)
}

func _CtiService_Join_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(JoinRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CtiServiceServer).Join(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.CtiService/Join",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CtiServiceServer).Join(ctx, req.(*JoinRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CtiService_Leave_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LeaveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CtiServiceServer).Leave(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.CtiService/Leave",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CtiServiceServer).Leave(ctx, req.(*LeaveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CtiService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CtiServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.CtiService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CtiServiceServer).List(ctx, req.(*ListRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CtiService_Info_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(InfoRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CtiServiceServer).Info(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.CtiService/Info",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CtiServiceServer).Info(ctx, req.(*InfoRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CtiService_ExtendsInstall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExtendsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CtiServiceServer).ExtendsInstall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.CtiService/ExtendsInstall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CtiServiceServer).ExtendsInstall(ctx, req.(*ExtendsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CtiService_ExtendsUpdate_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExtendsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CtiServiceServer).ExtendsUpdate(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.CtiService/ExtendsUpdate",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CtiServiceServer).ExtendsUpdate(ctx, req.(*ExtendsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _CtiService_ExtendsUninstall_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExtendsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(CtiServiceServer).ExtendsUninstall(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/service.CtiService/ExtendsUninstall",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(CtiServiceServer).ExtendsUninstall(ctx, req.(*ExtendsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// CtiService_ServiceDesc is the grpc.ServiceDesc for CtiService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var CtiService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.CtiService",
	HandlerType: (*CtiServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Join",
			Handler:    _CtiService_Join_Handler,
		},
		{
			MethodName: "Leave",
			Handler:    _CtiService_Leave_Handler,
		},
		{
			MethodName: "List",
			Handler:    _CtiService_List_Handler,
		},
		{
			MethodName: "Info",
			Handler:    _CtiService_Info_Handler,
		},
		{
			MethodName: "ExtendsInstall",
			Handler:    _CtiService_ExtendsInstall_Handler,
		},
		{
			MethodName: "ExtendsUpdate",
			Handler:    _CtiService_ExtendsUpdate_Handler,
		},
		{
			MethodName: "ExtendsUninstall",
			Handler:    _CtiService_ExtendsUninstall_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "ctl.proto",
}
