// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.19.6
// source: protos/service.proto

package protos

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

// HandlerClient is the client API for Handler service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type HandlerClient interface {
	// Sends a greeting
	Handle(ctx context.Context, in *HandleRequest, opts ...grpc.CallOption) (*HandleReply, error)
}

type handlerClient struct {
	cc grpc.ClientConnInterface
}

func NewHandlerClient(cc grpc.ClientConnInterface) HandlerClient {
	return &handlerClient{cc}
}

func (c *handlerClient) Handle(ctx context.Context, in *HandleRequest, opts ...grpc.CallOption) (*HandleReply, error) {
	out := new(HandleReply)
	err := c.cc.Invoke(ctx, "/Handler/Handle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// HandlerServer is the server API for Handler service.
// All implementations must embed UnimplementedHandlerServer
// for forward compatibility
type HandlerServer interface {
	// Sends a greeting
	Handle(context.Context, *HandleRequest) (*HandleReply, error)
	mustEmbedUnimplementedHandlerServer()
}

// UnimplementedHandlerServer must be embedded to have forward compatible implementations.
type UnimplementedHandlerServer struct {
}

func (UnimplementedHandlerServer) Handle(context.Context, *HandleRequest) (*HandleReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Handle not implemented")
}
func (UnimplementedHandlerServer) mustEmbedUnimplementedHandlerServer() {}

// UnsafeHandlerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to HandlerServer will
// result in compilation errors.
type UnsafeHandlerServer interface {
	mustEmbedUnimplementedHandlerServer()
}

func RegisterHandlerServer(s grpc.ServiceRegistrar, srv HandlerServer) {
	s.RegisterService(&Handler_ServiceDesc, srv)
}

func _Handler_Handle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HandleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(HandlerServer).Handle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/Handler/Handle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(HandlerServer).Handle(ctx, req.(*HandleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// Handler_ServiceDesc is the grpc.ServiceDesc for Handler service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var Handler_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "Handler",
	HandlerType: (*HandlerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Handle",
			Handler:    _Handler_Handle_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "protos/service.proto",
}