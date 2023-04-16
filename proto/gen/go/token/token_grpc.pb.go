// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.8
// source: token.proto

package token

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

// TokenSvcClient is the client API for TokenSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type TokenSvcClient interface {
	Register(ctx context.Context, in *HostInfo, opts ...grpc.CallOption) (*Token, error)
}

type tokenSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewTokenSvcClient(cc grpc.ClientConnInterface) TokenSvcClient {
	return &tokenSvcClient{cc}
}

func (c *tokenSvcClient) Register(ctx context.Context, in *HostInfo, opts ...grpc.CallOption) (*Token, error) {
	out := new(Token)
	err := c.cc.Invoke(ctx, "/TokenSvc/Register", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// TokenSvcServer is the server API for TokenSvc service.
// All implementations must embed UnimplementedTokenSvcServer
// for forward compatibility
type TokenSvcServer interface {
	Register(context.Context, *HostInfo) (*Token, error)
	mustEmbedUnimplementedTokenSvcServer()
}

// UnimplementedTokenSvcServer must be embedded to have forward compatible implementations.
type UnimplementedTokenSvcServer struct {
}

func (UnimplementedTokenSvcServer) Register(context.Context, *HostInfo) (*Token, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Register not implemented")
}
func (UnimplementedTokenSvcServer) mustEmbedUnimplementedTokenSvcServer() {}

// UnsafeTokenSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to TokenSvcServer will
// result in compilation errors.
type UnsafeTokenSvcServer interface {
	mustEmbedUnimplementedTokenSvcServer()
}

func RegisterTokenSvcServer(s grpc.ServiceRegistrar, srv TokenSvcServer) {
	s.RegisterService(&TokenSvc_ServiceDesc, srv)
}

func _TokenSvc_Register_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HostInfo)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(TokenSvcServer).Register(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/TokenSvc/Register",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(TokenSvcServer).Register(ctx, req.(*HostInfo))
	}
	return interceptor(ctx, in, info, handler)
}

// TokenSvc_ServiceDesc is the grpc.ServiceDesc for TokenSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var TokenSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "TokenSvc",
	HandlerType: (*TokenSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Register",
			Handler:    _TokenSvc_Register_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "token.proto",
}