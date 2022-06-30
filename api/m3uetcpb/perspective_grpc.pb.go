// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.12.4
// source: api/m3uetcpb/perspective.proto

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

// PerspectiveSvcClient is the client API for PerspectiveSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PerspectiveSvcClient interface {
	GetActivePerspective(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetActivePerspectiveResponse, error)
	SetActivePerspective(ctx context.Context, in *SetActivePerspectiveRequest, opts ...grpc.CallOption) (*Empty, error)
	SubscribeToPerspective(ctx context.Context, in *Empty, opts ...grpc.CallOption) (PerspectiveSvc_SubscribeToPerspectiveClient, error)
	UnsubscribeFromPerspective(ctx context.Context, in *UnsubscribeFromPerspectiveRequest, opts ...grpc.CallOption) (*Empty, error)
}

type perspectiveSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewPerspectiveSvcClient(cc grpc.ClientConnInterface) PerspectiveSvcClient {
	return &perspectiveSvcClient{cc}
}

func (c *perspectiveSvcClient) GetActivePerspective(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetActivePerspectiveResponse, error) {
	out := new(GetActivePerspectiveResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PerspectiveSvc/GetActivePerspective", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *perspectiveSvcClient) SetActivePerspective(ctx context.Context, in *SetActivePerspectiveRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PerspectiveSvc/SetActivePerspective", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *perspectiveSvcClient) SubscribeToPerspective(ctx context.Context, in *Empty, opts ...grpc.CallOption) (PerspectiveSvc_SubscribeToPerspectiveClient, error) {
	stream, err := c.cc.NewStream(ctx, &PerspectiveSvc_ServiceDesc.Streams[0], "/m3uetcpb.PerspectiveSvc/SubscribeToPerspective", opts...)
	if err != nil {
		return nil, err
	}
	x := &perspectiveSvcSubscribeToPerspectiveClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PerspectiveSvc_SubscribeToPerspectiveClient interface {
	Recv() (*SubscribeToPerspectiveResponse, error)
	grpc.ClientStream
}

type perspectiveSvcSubscribeToPerspectiveClient struct {
	grpc.ClientStream
}

func (x *perspectiveSvcSubscribeToPerspectiveClient) Recv() (*SubscribeToPerspectiveResponse, error) {
	m := new(SubscribeToPerspectiveResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *perspectiveSvcClient) UnsubscribeFromPerspective(ctx context.Context, in *UnsubscribeFromPerspectiveRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PerspectiveSvc/UnsubscribeFromPerspective", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PerspectiveSvcServer is the server API for PerspectiveSvc service.
// All implementations must embed UnimplementedPerspectiveSvcServer
// for forward compatibility
type PerspectiveSvcServer interface {
	GetActivePerspective(context.Context, *Empty) (*GetActivePerspectiveResponse, error)
	SetActivePerspective(context.Context, *SetActivePerspectiveRequest) (*Empty, error)
	SubscribeToPerspective(*Empty, PerspectiveSvc_SubscribeToPerspectiveServer) error
	UnsubscribeFromPerspective(context.Context, *UnsubscribeFromPerspectiveRequest) (*Empty, error)
	mustEmbedUnimplementedPerspectiveSvcServer()
}

// UnimplementedPerspectiveSvcServer must be embedded to have forward compatible implementations.
type UnimplementedPerspectiveSvcServer struct {
}

func (UnimplementedPerspectiveSvcServer) GetActivePerspective(context.Context, *Empty) (*GetActivePerspectiveResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetActivePerspective not implemented")
}
func (UnimplementedPerspectiveSvcServer) SetActivePerspective(context.Context, *SetActivePerspectiveRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SetActivePerspective not implemented")
}
func (UnimplementedPerspectiveSvcServer) SubscribeToPerspective(*Empty, PerspectiveSvc_SubscribeToPerspectiveServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeToPerspective not implemented")
}
func (UnimplementedPerspectiveSvcServer) UnsubscribeFromPerspective(context.Context, *UnsubscribeFromPerspectiveRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnsubscribeFromPerspective not implemented")
}
func (UnimplementedPerspectiveSvcServer) mustEmbedUnimplementedPerspectiveSvcServer() {}

// UnsafePerspectiveSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PerspectiveSvcServer will
// result in compilation errors.
type UnsafePerspectiveSvcServer interface {
	mustEmbedUnimplementedPerspectiveSvcServer()
}

func RegisterPerspectiveSvcServer(s grpc.ServiceRegistrar, srv PerspectiveSvcServer) {
	s.RegisterService(&PerspectiveSvc_ServiceDesc, srv)
}

func _PerspectiveSvc_GetActivePerspective_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PerspectiveSvcServer).GetActivePerspective(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PerspectiveSvc/GetActivePerspective",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PerspectiveSvcServer).GetActivePerspective(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PerspectiveSvc_SetActivePerspective_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(SetActivePerspectiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PerspectiveSvcServer).SetActivePerspective(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PerspectiveSvc/SetActivePerspective",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PerspectiveSvcServer).SetActivePerspective(ctx, req.(*SetActivePerspectiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PerspectiveSvc_SubscribeToPerspective_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PerspectiveSvcServer).SubscribeToPerspective(m, &perspectiveSvcSubscribeToPerspectiveServer{stream})
}

type PerspectiveSvc_SubscribeToPerspectiveServer interface {
	Send(*SubscribeToPerspectiveResponse) error
	grpc.ServerStream
}

type perspectiveSvcSubscribeToPerspectiveServer struct {
	grpc.ServerStream
}

func (x *perspectiveSvcSubscribeToPerspectiveServer) Send(m *SubscribeToPerspectiveResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _PerspectiveSvc_UnsubscribeFromPerspective_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnsubscribeFromPerspectiveRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PerspectiveSvcServer).UnsubscribeFromPerspective(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PerspectiveSvc/UnsubscribeFromPerspective",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PerspectiveSvcServer).UnsubscribeFromPerspective(ctx, req.(*UnsubscribeFromPerspectiveRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PerspectiveSvc_ServiceDesc is the grpc.ServiceDesc for PerspectiveSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PerspectiveSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "m3uetcpb.PerspectiveSvc",
	HandlerType: (*PerspectiveSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetActivePerspective",
			Handler:    _PerspectiveSvc_GetActivePerspective_Handler,
		},
		{
			MethodName: "SetActivePerspective",
			Handler:    _PerspectiveSvc_SetActivePerspective_Handler,
		},
		{
			MethodName: "UnsubscribeFromPerspective",
			Handler:    _PerspectiveSvc_UnsubscribeFromPerspective_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeToPerspective",
			Handler:       _PerspectiveSvc_SubscribeToPerspective_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/m3uetcpb/perspective.proto",
}
