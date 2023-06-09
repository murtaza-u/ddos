// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: admin.proto

package admin

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

// AdminSvcClient is the client API for AdminSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AdminSvcClient interface {
	GetIds(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Ids, error)
	Status(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Byts, error)
	DDos(ctx context.Context, in *Params, opts ...grpc.CallOption) (*Void, error)
}

type adminSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewAdminSvcClient(cc grpc.ClientConnInterface) AdminSvcClient {
	return &adminSvcClient{cc}
}

func (c *adminSvcClient) GetIds(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Ids, error) {
	out := new(Ids)
	err := c.cc.Invoke(ctx, "/AdminSvc/GetIds", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminSvcClient) Status(ctx context.Context, in *Void, opts ...grpc.CallOption) (*Byts, error) {
	out := new(Byts)
	err := c.cc.Invoke(ctx, "/AdminSvc/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *adminSvcClient) DDos(ctx context.Context, in *Params, opts ...grpc.CallOption) (*Void, error) {
	out := new(Void)
	err := c.cc.Invoke(ctx, "/AdminSvc/DDos", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AdminSvcServer is the server API for AdminSvc service.
// All implementations must embed UnimplementedAdminSvcServer
// for forward compatibility
type AdminSvcServer interface {
	GetIds(context.Context, *Void) (*Ids, error)
	Status(context.Context, *Void) (*Byts, error)
	DDos(context.Context, *Params) (*Void, error)
	mustEmbedUnimplementedAdminSvcServer()
}

// UnimplementedAdminSvcServer must be embedded to have forward compatible implementations.
type UnimplementedAdminSvcServer struct {
}

func (UnimplementedAdminSvcServer) GetIds(context.Context, *Void) (*Ids, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetIds not implemented")
}
func (UnimplementedAdminSvcServer) Status(context.Context, *Void) (*Byts, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedAdminSvcServer) DDos(context.Context, *Params) (*Void, error) {
	return nil, status.Errorf(codes.Unimplemented, "method DDos not implemented")
}
func (UnimplementedAdminSvcServer) mustEmbedUnimplementedAdminSvcServer() {}

// UnsafeAdminSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AdminSvcServer will
// result in compilation errors.
type UnsafeAdminSvcServer interface {
	mustEmbedUnimplementedAdminSvcServer()
}

func RegisterAdminSvcServer(s grpc.ServiceRegistrar, srv AdminSvcServer) {
	s.RegisterService(&AdminSvc_ServiceDesc, srv)
}

func _AdminSvc_GetIds_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminSvcServer).GetIds(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminSvc/GetIds",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminSvcServer).GetIds(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminSvc_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Void)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminSvcServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminSvc/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminSvcServer).Status(ctx, req.(*Void))
	}
	return interceptor(ctx, in, info, handler)
}

func _AdminSvc_DDos_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Params)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AdminSvcServer).DDos(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/AdminSvc/DDos",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AdminSvcServer).DDos(ctx, req.(*Params))
	}
	return interceptor(ctx, in, info, handler)
}

// AdminSvc_ServiceDesc is the grpc.ServiceDesc for AdminSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AdminSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "AdminSvc",
	HandlerType: (*AdminSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetIds",
			Handler:    _AdminSvc_GetIds_Handler,
		},
		{
			MethodName: "Status",
			Handler:    _AdminSvc_Status_Handler,
		},
		{
			MethodName: "DDos",
			Handler:    _AdminSvc_DDos_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "admin.proto",
}
