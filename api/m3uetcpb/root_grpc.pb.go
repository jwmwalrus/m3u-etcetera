// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: api/m3uetcpb/root.proto

package m3uetcpb

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

// RootSvcClient is the client API for RootSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type RootSvcClient interface {
	Status(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*StatusResponse, error)
	Off(ctx context.Context, in *OffRequest, opts ...grpc.CallOption) (*OffResponse, error)
}

type rootSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewRootSvcClient(cc grpc.ClientConnInterface) RootSvcClient {
	return &rootSvcClient{cc}
}

func (c *rootSvcClient) Status(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*StatusResponse, error) {
	out := new(StatusResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.RootSvc/Status", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *rootSvcClient) Off(ctx context.Context, in *OffRequest, opts ...grpc.CallOption) (*OffResponse, error) {
	out := new(OffResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.RootSvc/Off", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// RootSvcServer is the server API for RootSvc service.
// All implementations must embed UnimplementedRootSvcServer
// for forward compatibility
type RootSvcServer interface {
	Status(context.Context, *Empty) (*StatusResponse, error)
	Off(context.Context, *OffRequest) (*OffResponse, error)
	mustEmbedUnimplementedRootSvcServer()
}

// UnimplementedRootSvcServer must be embedded to have forward compatible implementations.
type UnimplementedRootSvcServer struct {
}

func (UnimplementedRootSvcServer) Status(context.Context, *Empty) (*StatusResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Status not implemented")
}
func (UnimplementedRootSvcServer) Off(context.Context, *OffRequest) (*OffResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Off not implemented")
}
func (UnimplementedRootSvcServer) mustEmbedUnimplementedRootSvcServer() {}

// UnsafeRootSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to RootSvcServer will
// result in compilation errors.
type UnsafeRootSvcServer interface {
	mustEmbedUnimplementedRootSvcServer()
}

func RegisterRootSvcServer(s grpc.ServiceRegistrar, srv RootSvcServer) {
	s.RegisterService(&RootSvc_ServiceDesc, srv)
}

func _RootSvc_Status_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RootSvcServer).Status(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.RootSvc/Status",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RootSvcServer).Status(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _RootSvc_Off_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(OffRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(RootSvcServer).Off(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.RootSvc/Off",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(RootSvcServer).Off(ctx, req.(*OffRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// RootSvc_ServiceDesc is the grpc.ServiceDesc for RootSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var RootSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "m3uetcpb.RootSvc",
	HandlerType: (*RootSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Status",
			Handler:    _RootSvc_Status_Handler,
		},
		{
			MethodName: "Off",
			Handler:    _RootSvc_Off_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "api/m3uetcpb/root.proto",
}
