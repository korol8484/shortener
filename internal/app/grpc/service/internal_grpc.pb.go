// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v3.12.4
// source: proto/service/internal.proto

package service

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	contract "github.com/korol8484/shortener/internal/app/grpc/service/contract"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	Internal_Ping_FullMethodName        = "/service.Internal/Ping"
	Internal_HandleShort_FullMethodName = "/service.Internal/HandleShort"
	Internal_HandleGet_FullMethodName   = "/service.Internal/HandleGet"
	Internal_HandleBatch_FullMethodName = "/service.Internal/HandleBatch"
	Internal_UserURL_FullMethodName     = "/service.Internal/UserURL"
	Internal_BatchDelete_FullMethodName = "/service.Internal/BatchDelete"
	Internal_Stats_FullMethodName       = "/service.Internal/Stats"
)

// InternalClient is the client API for Internal service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type InternalClient interface {
	Ping(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
	HandleShort(ctx context.Context, in *contract.RequestShort, opts ...grpc.CallOption) (*contract.ResponseShort, error)
	HandleGet(ctx context.Context, in *contract.RequestFindByAlias, opts ...grpc.CallOption) (*contract.ResponseShort, error)
	HandleBatch(ctx context.Context, in *contract.RequestBatch, opts ...grpc.CallOption) (*contract.ResponseBatch, error)
	UserURL(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*contract.ResponseUserUrl, error)
	BatchDelete(ctx context.Context, in *contract.RequestBatchDelete, opts ...grpc.CallOption) (*empty.Empty, error)
	Stats(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*contract.ResponseStats, error)
}

type internalClient struct {
	cc grpc.ClientConnInterface
}

func NewInternalClient(cc grpc.ClientConnInterface) InternalClient {
	return &internalClient{cc}
}

func (c *internalClient) Ping(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Internal_Ping_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *internalClient) HandleShort(ctx context.Context, in *contract.RequestShort, opts ...grpc.CallOption) (*contract.ResponseShort, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(contract.ResponseShort)
	err := c.cc.Invoke(ctx, Internal_HandleShort_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *internalClient) HandleGet(ctx context.Context, in *contract.RequestFindByAlias, opts ...grpc.CallOption) (*contract.ResponseShort, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(contract.ResponseShort)
	err := c.cc.Invoke(ctx, Internal_HandleGet_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *internalClient) HandleBatch(ctx context.Context, in *contract.RequestBatch, opts ...grpc.CallOption) (*contract.ResponseBatch, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(contract.ResponseBatch)
	err := c.cc.Invoke(ctx, Internal_HandleBatch_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *internalClient) UserURL(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*contract.ResponseUserUrl, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(contract.ResponseUserUrl)
	err := c.cc.Invoke(ctx, Internal_UserURL_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *internalClient) BatchDelete(ctx context.Context, in *contract.RequestBatchDelete, opts ...grpc.CallOption) (*empty.Empty, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, Internal_BatchDelete_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *internalClient) Stats(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*contract.ResponseStats, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(contract.ResponseStats)
	err := c.cc.Invoke(ctx, Internal_Stats_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// InternalServer is the server API for Internal service.
// All implementations must embed UnimplementedInternalServer
// for forward compatibility.
type InternalServer interface {
	Ping(context.Context, *empty.Empty) (*empty.Empty, error)
	HandleShort(context.Context, *contract.RequestShort) (*contract.ResponseShort, error)
	HandleGet(context.Context, *contract.RequestFindByAlias) (*contract.ResponseShort, error)
	HandleBatch(context.Context, *contract.RequestBatch) (*contract.ResponseBatch, error)
	UserURL(context.Context, *empty.Empty) (*contract.ResponseUserUrl, error)
	BatchDelete(context.Context, *contract.RequestBatchDelete) (*empty.Empty, error)
	Stats(context.Context, *empty.Empty) (*contract.ResponseStats, error)
	mustEmbedUnimplementedInternalServer()
}

// UnimplementedInternalServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedInternalServer struct{}

func (UnimplementedInternalServer) Ping(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedInternalServer) HandleShort(context.Context, *contract.RequestShort) (*contract.ResponseShort, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleShort not implemented")
}
func (UnimplementedInternalServer) HandleGet(context.Context, *contract.RequestFindByAlias) (*contract.ResponseShort, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleGet not implemented")
}
func (UnimplementedInternalServer) HandleBatch(context.Context, *contract.RequestBatch) (*contract.ResponseBatch, error) {
	return nil, status.Errorf(codes.Unimplemented, "method HandleBatch not implemented")
}
func (UnimplementedInternalServer) UserURL(context.Context, *empty.Empty) (*contract.ResponseUserUrl, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UserURL not implemented")
}
func (UnimplementedInternalServer) BatchDelete(context.Context, *contract.RequestBatchDelete) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method BatchDelete not implemented")
}
func (UnimplementedInternalServer) Stats(context.Context, *empty.Empty) (*contract.ResponseStats, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Stats not implemented")
}
func (UnimplementedInternalServer) mustEmbedUnimplementedInternalServer() {}
func (UnimplementedInternalServer) testEmbeddedByValue()                  {}

// UnsafeInternalServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to InternalServer will
// result in compilation errors.
type UnsafeInternalServer interface {
	mustEmbedUnimplementedInternalServer()
}

func RegisterInternalServer(s grpc.ServiceRegistrar, srv InternalServer) {
	// If the following call pancis, it indicates UnimplementedInternalServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&Internal_ServiceDesc, srv)
}

func _Internal_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternalServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Internal_Ping_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternalServer).Ping(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Internal_HandleShort_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(contract.RequestShort)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternalServer).HandleShort(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Internal_HandleShort_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternalServer).HandleShort(ctx, req.(*contract.RequestShort))
	}
	return interceptor(ctx, in, info, handler)
}

func _Internal_HandleGet_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(contract.RequestFindByAlias)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternalServer).HandleGet(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Internal_HandleGet_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternalServer).HandleGet(ctx, req.(*contract.RequestFindByAlias))
	}
	return interceptor(ctx, in, info, handler)
}

func _Internal_HandleBatch_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(contract.RequestBatch)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternalServer).HandleBatch(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Internal_HandleBatch_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternalServer).HandleBatch(ctx, req.(*contract.RequestBatch))
	}
	return interceptor(ctx, in, info, handler)
}

func _Internal_UserURL_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternalServer).UserURL(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Internal_UserURL_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternalServer).UserURL(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _Internal_BatchDelete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(contract.RequestBatchDelete)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternalServer).BatchDelete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Internal_BatchDelete_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternalServer).BatchDelete(ctx, req.(*contract.RequestBatchDelete))
	}
	return interceptor(ctx, in, info, handler)
}

func _Internal_Stats_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(InternalServer).Stats(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: Internal_Stats_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(InternalServer).Stats(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// Internal_ServiceDesc is the grpc.ServiceDesc for Internal service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Internal_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "service.Internal",
	HandlerType: (*InternalServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Ping",
			Handler:    _Internal_Ping_Handler,
		},
		{
			MethodName: "HandleShort",
			Handler:    _Internal_HandleShort_Handler,
		},
		{
			MethodName: "HandleGet",
			Handler:    _Internal_HandleGet_Handler,
		},
		{
			MethodName: "HandleBatch",
			Handler:    _Internal_HandleBatch_Handler,
		},
		{
			MethodName: "UserURL",
			Handler:    _Internal_UserURL_Handler,
		},
		{
			MethodName: "BatchDelete",
			Handler:    _Internal_BatchDelete_Handler,
		},
		{
			MethodName: "Stats",
			Handler:    _Internal_Stats_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/service/internal.proto",
}