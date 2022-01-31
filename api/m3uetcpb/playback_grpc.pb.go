// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

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

// PlaybackSvcClient is the client API for PlaybackSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlaybackSvcClient interface {
	GetPlayback(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetPlaybackResponse, error)
	GetPlaybackList(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetPlaybackListResponse, error)
	ExecutePlaybackAction(ctx context.Context, in *ExecutePlaybackActionRequest, opts ...grpc.CallOption) (*Empty, error)
	SubscribeToPlayback(ctx context.Context, in *Empty, opts ...grpc.CallOption) (PlaybackSvc_SubscribeToPlaybackClient, error)
	UnsubscribeFromPlayback(ctx context.Context, in *UnsubscribeFromPlaybackRequest, opts ...grpc.CallOption) (*Empty, error)
}

type playbackSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewPlaybackSvcClient(cc grpc.ClientConnInterface) PlaybackSvcClient {
	return &playbackSvcClient{cc}
}

func (c *playbackSvcClient) GetPlayback(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetPlaybackResponse, error) {
	out := new(GetPlaybackResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybackSvc/GetPlayback", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playbackSvcClient) GetPlaybackList(ctx context.Context, in *Empty, opts ...grpc.CallOption) (*GetPlaybackListResponse, error) {
	out := new(GetPlaybackListResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybackSvc/GetPlaybackList", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playbackSvcClient) ExecutePlaybackAction(ctx context.Context, in *ExecutePlaybackActionRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybackSvc/ExecutePlaybackAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playbackSvcClient) SubscribeToPlayback(ctx context.Context, in *Empty, opts ...grpc.CallOption) (PlaybackSvc_SubscribeToPlaybackClient, error) {
	stream, err := c.cc.NewStream(ctx, &PlaybackSvc_ServiceDesc.Streams[0], "/m3uetcpb.PlaybackSvc/SubscribeToPlayback", opts...)
	if err != nil {
		return nil, err
	}
	x := &playbackSvcSubscribeToPlaybackClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PlaybackSvc_SubscribeToPlaybackClient interface {
	Recv() (*SubscribeToPlaybackResponse, error)
	grpc.ClientStream
}

type playbackSvcSubscribeToPlaybackClient struct {
	grpc.ClientStream
}

func (x *playbackSvcSubscribeToPlaybackClient) Recv() (*SubscribeToPlaybackResponse, error) {
	m := new(SubscribeToPlaybackResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *playbackSvcClient) UnsubscribeFromPlayback(ctx context.Context, in *UnsubscribeFromPlaybackRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybackSvc/UnsubscribeFromPlayback", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlaybackSvcServer is the server API for PlaybackSvc service.
// All implementations must embed UnimplementedPlaybackSvcServer
// for forward compatibility
type PlaybackSvcServer interface {
	GetPlayback(context.Context, *Empty) (*GetPlaybackResponse, error)
	GetPlaybackList(context.Context, *Empty) (*GetPlaybackListResponse, error)
	ExecutePlaybackAction(context.Context, *ExecutePlaybackActionRequest) (*Empty, error)
	SubscribeToPlayback(*Empty, PlaybackSvc_SubscribeToPlaybackServer) error
	UnsubscribeFromPlayback(context.Context, *UnsubscribeFromPlaybackRequest) (*Empty, error)
	mustEmbedUnimplementedPlaybackSvcServer()
}

// UnimplementedPlaybackSvcServer must be embedded to have forward compatible implementations.
type UnimplementedPlaybackSvcServer struct {
}

func (UnimplementedPlaybackSvcServer) GetPlayback(context.Context, *Empty) (*GetPlaybackResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlayback not implemented")
}
func (UnimplementedPlaybackSvcServer) GetPlaybackList(context.Context, *Empty) (*GetPlaybackListResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlaybackList not implemented")
}
func (UnimplementedPlaybackSvcServer) ExecutePlaybackAction(context.Context, *ExecutePlaybackActionRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecutePlaybackAction not implemented")
}
func (UnimplementedPlaybackSvcServer) SubscribeToPlayback(*Empty, PlaybackSvc_SubscribeToPlaybackServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeToPlayback not implemented")
}
func (UnimplementedPlaybackSvcServer) UnsubscribeFromPlayback(context.Context, *UnsubscribeFromPlaybackRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnsubscribeFromPlayback not implemented")
}
func (UnimplementedPlaybackSvcServer) mustEmbedUnimplementedPlaybackSvcServer() {}

// UnsafePlaybackSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlaybackSvcServer will
// result in compilation errors.
type UnsafePlaybackSvcServer interface {
	mustEmbedUnimplementedPlaybackSvcServer()
}

func RegisterPlaybackSvcServer(s grpc.ServiceRegistrar, srv PlaybackSvcServer) {
	s.RegisterService(&PlaybackSvc_ServiceDesc, srv)
}

func _PlaybackSvc_GetPlayback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybackSvcServer).GetPlayback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybackSvc/GetPlayback",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybackSvcServer).GetPlayback(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaybackSvc_GetPlaybackList_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybackSvcServer).GetPlaybackList(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybackSvc/GetPlaybackList",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybackSvcServer).GetPlaybackList(ctx, req.(*Empty))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaybackSvc_ExecutePlaybackAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecutePlaybackActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybackSvcServer).ExecutePlaybackAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybackSvc/ExecutePlaybackAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybackSvcServer).ExecutePlaybackAction(ctx, req.(*ExecutePlaybackActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaybackSvc_SubscribeToPlayback_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PlaybackSvcServer).SubscribeToPlayback(m, &playbackSvcSubscribeToPlaybackServer{stream})
}

type PlaybackSvc_SubscribeToPlaybackServer interface {
	Send(*SubscribeToPlaybackResponse) error
	grpc.ServerStream
}

type playbackSvcSubscribeToPlaybackServer struct {
	grpc.ServerStream
}

func (x *playbackSvcSubscribeToPlaybackServer) Send(m *SubscribeToPlaybackResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _PlaybackSvc_UnsubscribeFromPlayback_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnsubscribeFromPlaybackRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybackSvcServer).UnsubscribeFromPlayback(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybackSvc/UnsubscribeFromPlayback",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybackSvcServer).UnsubscribeFromPlayback(ctx, req.(*UnsubscribeFromPlaybackRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PlaybackSvc_ServiceDesc is the grpc.ServiceDesc for PlaybackSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PlaybackSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "m3uetcpb.PlaybackSvc",
	HandlerType: (*PlaybackSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPlayback",
			Handler:    _PlaybackSvc_GetPlayback_Handler,
		},
		{
			MethodName: "GetPlaybackList",
			Handler:    _PlaybackSvc_GetPlaybackList_Handler,
		},
		{
			MethodName: "ExecutePlaybackAction",
			Handler:    _PlaybackSvc_ExecutePlaybackAction_Handler,
		},
		{
			MethodName: "UnsubscribeFromPlayback",
			Handler:    _PlaybackSvc_UnsubscribeFromPlayback_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeToPlayback",
			Handler:       _PlaybackSvc_SubscribeToPlayback_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/m3uetcpb/playback.proto",
}
