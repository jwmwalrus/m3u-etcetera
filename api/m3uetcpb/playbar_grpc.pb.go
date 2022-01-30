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

// PlaybarSvcClient is the client API for PlaybarSvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PlaybarSvcClient interface {
	GetPlaybar(ctx context.Context, in *GetPlaybarRequest, opts ...grpc.CallOption) (*GetPlaybarResponse, error)
	GetPlaylist(ctx context.Context, in *GetPlaylistRequest, opts ...grpc.CallOption) (*GetPlaylistResponse, error)
	GetAllPlaylists(ctx context.Context, in *GetAllPlaylistsRequest, opts ...grpc.CallOption) (*GetAllPlaylistsResponse, error)
	GetPlaylistGroup(ctx context.Context, in *GetPlaylistGroupRequest, opts ...grpc.CallOption) (*GetPlaylistGroupResponse, error)
	GetAllPlaylistGroups(ctx context.Context, in *GetAllPlaylistGroupsRequest, opts ...grpc.CallOption) (*GetAllPlaylistGroupsResponse, error)
	ExecutePlaybarAction(ctx context.Context, in *ExecutePlaybarActionRequest, opts ...grpc.CallOption) (*Empty, error)
	ExecutePlaylistAction(ctx context.Context, in *ExecutePlaylistActionRequest, opts ...grpc.CallOption) (*ExecutePlaylistActionResponse, error)
	ExecutePlaylistGroupAction(ctx context.Context, in *ExecutePlaylistGroupActionRequest, opts ...grpc.CallOption) (*ExecutePlaylistGroupActionResponse, error)
	ExecutePlaylistTrackAction(ctx context.Context, in *ExecutePlaylistTrackActionRequest, opts ...grpc.CallOption) (*Empty, error)
	ImportPlaylists(ctx context.Context, in *ImportPlaylistsRequest, opts ...grpc.CallOption) (PlaybarSvc_ImportPlaylistsClient, error)
	SubscribeToPlaybarStore(ctx context.Context, in *Empty, opts ...grpc.CallOption) (PlaybarSvc_SubscribeToPlaybarStoreClient, error)
	UnsubscribeFromPlaybarStore(ctx context.Context, in *UnsubscribeFromPlaybarStoreRequest, opts ...grpc.CallOption) (*Empty, error)
}

type playbarSvcClient struct {
	cc grpc.ClientConnInterface
}

func NewPlaybarSvcClient(cc grpc.ClientConnInterface) PlaybarSvcClient {
	return &playbarSvcClient{cc}
}

func (c *playbarSvcClient) GetPlaybar(ctx context.Context, in *GetPlaybarRequest, opts ...grpc.CallOption) (*GetPlaybarResponse, error) {
	out := new(GetPlaybarResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybarSvc/GetPlaybar", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playbarSvcClient) GetPlaylist(ctx context.Context, in *GetPlaylistRequest, opts ...grpc.CallOption) (*GetPlaylistResponse, error) {
	out := new(GetPlaylistResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybarSvc/GetPlaylist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playbarSvcClient) GetAllPlaylists(ctx context.Context, in *GetAllPlaylistsRequest, opts ...grpc.CallOption) (*GetAllPlaylistsResponse, error) {
	out := new(GetAllPlaylistsResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybarSvc/GetAllPlaylists", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playbarSvcClient) GetPlaylistGroup(ctx context.Context, in *GetPlaylistGroupRequest, opts ...grpc.CallOption) (*GetPlaylistGroupResponse, error) {
	out := new(GetPlaylistGroupResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybarSvc/GetPlaylistGroup", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playbarSvcClient) GetAllPlaylistGroups(ctx context.Context, in *GetAllPlaylistGroupsRequest, opts ...grpc.CallOption) (*GetAllPlaylistGroupsResponse, error) {
	out := new(GetAllPlaylistGroupsResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybarSvc/GetAllPlaylistGroups", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playbarSvcClient) ExecutePlaybarAction(ctx context.Context, in *ExecutePlaybarActionRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybarSvc/ExecutePlaybarAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playbarSvcClient) ExecutePlaylistAction(ctx context.Context, in *ExecutePlaylistActionRequest, opts ...grpc.CallOption) (*ExecutePlaylistActionResponse, error) {
	out := new(ExecutePlaylistActionResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybarSvc/ExecutePlaylistAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playbarSvcClient) ExecutePlaylistGroupAction(ctx context.Context, in *ExecutePlaylistGroupActionRequest, opts ...grpc.CallOption) (*ExecutePlaylistGroupActionResponse, error) {
	out := new(ExecutePlaylistGroupActionResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybarSvc/ExecutePlaylistGroupAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playbarSvcClient) ExecutePlaylistTrackAction(ctx context.Context, in *ExecutePlaylistTrackActionRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybarSvc/ExecutePlaylistTrackAction", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *playbarSvcClient) ImportPlaylists(ctx context.Context, in *ImportPlaylistsRequest, opts ...grpc.CallOption) (PlaybarSvc_ImportPlaylistsClient, error) {
	stream, err := c.cc.NewStream(ctx, &PlaybarSvc_ServiceDesc.Streams[0], "/m3uetcpb.PlaybarSvc/ImportPlaylists", opts...)
	if err != nil {
		return nil, err
	}
	x := &playbarSvcImportPlaylistsClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PlaybarSvc_ImportPlaylistsClient interface {
	Recv() (*ImportPlaylistsResponse, error)
	grpc.ClientStream
}

type playbarSvcImportPlaylistsClient struct {
	grpc.ClientStream
}

func (x *playbarSvcImportPlaylistsClient) Recv() (*ImportPlaylistsResponse, error) {
	m := new(ImportPlaylistsResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *playbarSvcClient) SubscribeToPlaybarStore(ctx context.Context, in *Empty, opts ...grpc.CallOption) (PlaybarSvc_SubscribeToPlaybarStoreClient, error) {
	stream, err := c.cc.NewStream(ctx, &PlaybarSvc_ServiceDesc.Streams[1], "/m3uetcpb.PlaybarSvc/SubscribeToPlaybarStore", opts...)
	if err != nil {
		return nil, err
	}
	x := &playbarSvcSubscribeToPlaybarStoreClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type PlaybarSvc_SubscribeToPlaybarStoreClient interface {
	Recv() (*SubscribeToPlaybarStoreResponse, error)
	grpc.ClientStream
}

type playbarSvcSubscribeToPlaybarStoreClient struct {
	grpc.ClientStream
}

func (x *playbarSvcSubscribeToPlaybarStoreClient) Recv() (*SubscribeToPlaybarStoreResponse, error) {
	m := new(SubscribeToPlaybarStoreResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *playbarSvcClient) UnsubscribeFromPlaybarStore(ctx context.Context, in *UnsubscribeFromPlaybarStoreRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/m3uetcpb.PlaybarSvc/UnsubscribeFromPlaybarStore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PlaybarSvcServer is the server API for PlaybarSvc service.
// All implementations must embed UnimplementedPlaybarSvcServer
// for forward compatibility
type PlaybarSvcServer interface {
	GetPlaybar(context.Context, *GetPlaybarRequest) (*GetPlaybarResponse, error)
	GetPlaylist(context.Context, *GetPlaylistRequest) (*GetPlaylistResponse, error)
	GetAllPlaylists(context.Context, *GetAllPlaylistsRequest) (*GetAllPlaylistsResponse, error)
	GetPlaylistGroup(context.Context, *GetPlaylistGroupRequest) (*GetPlaylistGroupResponse, error)
	GetAllPlaylistGroups(context.Context, *GetAllPlaylistGroupsRequest) (*GetAllPlaylistGroupsResponse, error)
	ExecutePlaybarAction(context.Context, *ExecutePlaybarActionRequest) (*Empty, error)
	ExecutePlaylistAction(context.Context, *ExecutePlaylistActionRequest) (*ExecutePlaylistActionResponse, error)
	ExecutePlaylistGroupAction(context.Context, *ExecutePlaylistGroupActionRequest) (*ExecutePlaylistGroupActionResponse, error)
	ExecutePlaylistTrackAction(context.Context, *ExecutePlaylistTrackActionRequest) (*Empty, error)
	ImportPlaylists(*ImportPlaylistsRequest, PlaybarSvc_ImportPlaylistsServer) error
	SubscribeToPlaybarStore(*Empty, PlaybarSvc_SubscribeToPlaybarStoreServer) error
	UnsubscribeFromPlaybarStore(context.Context, *UnsubscribeFromPlaybarStoreRequest) (*Empty, error)
	mustEmbedUnimplementedPlaybarSvcServer()
}

// UnimplementedPlaybarSvcServer must be embedded to have forward compatible implementations.
type UnimplementedPlaybarSvcServer struct {
}

func (UnimplementedPlaybarSvcServer) GetPlaybar(context.Context, *GetPlaybarRequest) (*GetPlaybarResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlaybar not implemented")
}
func (UnimplementedPlaybarSvcServer) GetPlaylist(context.Context, *GetPlaylistRequest) (*GetPlaylistResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlaylist not implemented")
}
func (UnimplementedPlaybarSvcServer) GetAllPlaylists(context.Context, *GetAllPlaylistsRequest) (*GetAllPlaylistsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllPlaylists not implemented")
}
func (UnimplementedPlaybarSvcServer) GetPlaylistGroup(context.Context, *GetPlaylistGroupRequest) (*GetPlaylistGroupResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPlaylistGroup not implemented")
}
func (UnimplementedPlaybarSvcServer) GetAllPlaylistGroups(context.Context, *GetAllPlaylistGroupsRequest) (*GetAllPlaylistGroupsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetAllPlaylistGroups not implemented")
}
func (UnimplementedPlaybarSvcServer) ExecutePlaybarAction(context.Context, *ExecutePlaybarActionRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecutePlaybarAction not implemented")
}
func (UnimplementedPlaybarSvcServer) ExecutePlaylistAction(context.Context, *ExecutePlaylistActionRequest) (*ExecutePlaylistActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecutePlaylistAction not implemented")
}
func (UnimplementedPlaybarSvcServer) ExecutePlaylistGroupAction(context.Context, *ExecutePlaylistGroupActionRequest) (*ExecutePlaylistGroupActionResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecutePlaylistGroupAction not implemented")
}
func (UnimplementedPlaybarSvcServer) ExecutePlaylistTrackAction(context.Context, *ExecutePlaylistTrackActionRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method ExecutePlaylistTrackAction not implemented")
}
func (UnimplementedPlaybarSvcServer) ImportPlaylists(*ImportPlaylistsRequest, PlaybarSvc_ImportPlaylistsServer) error {
	return status.Errorf(codes.Unimplemented, "method ImportPlaylists not implemented")
}
func (UnimplementedPlaybarSvcServer) SubscribeToPlaybarStore(*Empty, PlaybarSvc_SubscribeToPlaybarStoreServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeToPlaybarStore not implemented")
}
func (UnimplementedPlaybarSvcServer) UnsubscribeFromPlaybarStore(context.Context, *UnsubscribeFromPlaybarStoreRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnsubscribeFromPlaybarStore not implemented")
}
func (UnimplementedPlaybarSvcServer) mustEmbedUnimplementedPlaybarSvcServer() {}

// UnsafePlaybarSvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PlaybarSvcServer will
// result in compilation errors.
type UnsafePlaybarSvcServer interface {
	mustEmbedUnimplementedPlaybarSvcServer()
}

func RegisterPlaybarSvcServer(s grpc.ServiceRegistrar, srv PlaybarSvcServer) {
	s.RegisterService(&PlaybarSvc_ServiceDesc, srv)
}

func _PlaybarSvc_GetPlaybar_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPlaybarRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybarSvcServer).GetPlaybar(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybarSvc/GetPlaybar",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybarSvcServer).GetPlaybar(ctx, req.(*GetPlaybarRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaybarSvc_GetPlaylist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPlaylistRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybarSvcServer).GetPlaylist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybarSvc/GetPlaylist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybarSvcServer).GetPlaylist(ctx, req.(*GetPlaylistRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaybarSvc_GetAllPlaylists_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllPlaylistsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybarSvcServer).GetAllPlaylists(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybarSvc/GetAllPlaylists",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybarSvcServer).GetAllPlaylists(ctx, req.(*GetAllPlaylistsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaybarSvc_GetPlaylistGroup_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPlaylistGroupRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybarSvcServer).GetPlaylistGroup(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybarSvc/GetPlaylistGroup",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybarSvcServer).GetPlaylistGroup(ctx, req.(*GetPlaylistGroupRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaybarSvc_GetAllPlaylistGroups_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAllPlaylistGroupsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybarSvcServer).GetAllPlaylistGroups(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybarSvc/GetAllPlaylistGroups",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybarSvcServer).GetAllPlaylistGroups(ctx, req.(*GetAllPlaylistGroupsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaybarSvc_ExecutePlaybarAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecutePlaybarActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybarSvcServer).ExecutePlaybarAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybarSvc/ExecutePlaybarAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybarSvcServer).ExecutePlaybarAction(ctx, req.(*ExecutePlaybarActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaybarSvc_ExecutePlaylistAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecutePlaylistActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybarSvcServer).ExecutePlaylistAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybarSvc/ExecutePlaylistAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybarSvcServer).ExecutePlaylistAction(ctx, req.(*ExecutePlaylistActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaybarSvc_ExecutePlaylistGroupAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecutePlaylistGroupActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybarSvcServer).ExecutePlaylistGroupAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybarSvc/ExecutePlaylistGroupAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybarSvcServer).ExecutePlaylistGroupAction(ctx, req.(*ExecutePlaylistGroupActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaybarSvc_ExecutePlaylistTrackAction_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ExecutePlaylistTrackActionRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybarSvcServer).ExecutePlaylistTrackAction(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybarSvc/ExecutePlaylistTrackAction",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybarSvcServer).ExecutePlaylistTrackAction(ctx, req.(*ExecutePlaylistTrackActionRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _PlaybarSvc_ImportPlaylists_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(ImportPlaylistsRequest)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PlaybarSvcServer).ImportPlaylists(m, &playbarSvcImportPlaylistsServer{stream})
}

type PlaybarSvc_ImportPlaylistsServer interface {
	Send(*ImportPlaylistsResponse) error
	grpc.ServerStream
}

type playbarSvcImportPlaylistsServer struct {
	grpc.ServerStream
}

func (x *playbarSvcImportPlaylistsServer) Send(m *ImportPlaylistsResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _PlaybarSvc_SubscribeToPlaybarStore_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(PlaybarSvcServer).SubscribeToPlaybarStore(m, &playbarSvcSubscribeToPlaybarStoreServer{stream})
}

type PlaybarSvc_SubscribeToPlaybarStoreServer interface {
	Send(*SubscribeToPlaybarStoreResponse) error
	grpc.ServerStream
}

type playbarSvcSubscribeToPlaybarStoreServer struct {
	grpc.ServerStream
}

func (x *playbarSvcSubscribeToPlaybarStoreServer) Send(m *SubscribeToPlaybarStoreResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _PlaybarSvc_UnsubscribeFromPlaybarStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnsubscribeFromPlaybarStoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PlaybarSvcServer).UnsubscribeFromPlaybarStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.PlaybarSvc/UnsubscribeFromPlaybarStore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PlaybarSvcServer).UnsubscribeFromPlaybarStore(ctx, req.(*UnsubscribeFromPlaybarStoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// PlaybarSvc_ServiceDesc is the grpc.ServiceDesc for PlaybarSvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PlaybarSvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "m3uetcpb.PlaybarSvc",
	HandlerType: (*PlaybarSvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetPlaybar",
			Handler:    _PlaybarSvc_GetPlaybar_Handler,
		},
		{
			MethodName: "GetPlaylist",
			Handler:    _PlaybarSvc_GetPlaylist_Handler,
		},
		{
			MethodName: "GetAllPlaylists",
			Handler:    _PlaybarSvc_GetAllPlaylists_Handler,
		},
		{
			MethodName: "GetPlaylistGroup",
			Handler:    _PlaybarSvc_GetPlaylistGroup_Handler,
		},
		{
			MethodName: "GetAllPlaylistGroups",
			Handler:    _PlaybarSvc_GetAllPlaylistGroups_Handler,
		},
		{
			MethodName: "ExecutePlaybarAction",
			Handler:    _PlaybarSvc_ExecutePlaybarAction_Handler,
		},
		{
			MethodName: "ExecutePlaylistAction",
			Handler:    _PlaybarSvc_ExecutePlaylistAction_Handler,
		},
		{
			MethodName: "ExecutePlaylistGroupAction",
			Handler:    _PlaybarSvc_ExecutePlaylistGroupAction_Handler,
		},
		{
			MethodName: "ExecutePlaylistTrackAction",
			Handler:    _PlaybarSvc_ExecutePlaylistTrackAction_Handler,
		},
		{
			MethodName: "UnsubscribeFromPlaybarStore",
			Handler:    _PlaybarSvc_UnsubscribeFromPlaybarStore_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "ImportPlaylists",
			Handler:       _PlaybarSvc_ImportPlaylists_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "SubscribeToPlaybarStore",
			Handler:       _PlaybarSvc_SubscribeToPlaybarStore_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/m3uetcpb/playbar.proto",
}
