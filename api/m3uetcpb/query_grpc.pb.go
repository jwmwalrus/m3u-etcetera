// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.21.12
// source: api/m3uetcpb/query.proto

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

// QuerySvcClient is the client API for QuerySvc service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type QuerySvcClient interface {
	GetQuery(ctx context.Context, in *GetQueryRequest, opts ...grpc.CallOption) (*GetQueryResponse, error)
	GetQueries(ctx context.Context, in *GetQueriesRequest, opts ...grpc.CallOption) (*GetQueriesResponse, error)
	AddQuery(ctx context.Context, in *AddQueryRequest, opts ...grpc.CallOption) (*AddQueryResponse, error)
	UpdateQuery(ctx context.Context, in *UpdateQueryRequest, opts ...grpc.CallOption) (*Empty, error)
	RemoveQuery(ctx context.Context, in *RemoveQueryRequest, opts ...grpc.CallOption) (*Empty, error)
	QueryBy(ctx context.Context, in *QueryByRequest, opts ...grpc.CallOption) (*QueryByResponse, error)
	QueryInPlaylist(ctx context.Context, in *QueryInPlaylistRequest, opts ...grpc.CallOption) (*QueryInPlaylistResponse, error)
	QueryInQueue(ctx context.Context, in *QueryInQueueRequest, opts ...grpc.CallOption) (*Empty, error)
	SubscribeToQueryStore(ctx context.Context, in *Empty, opts ...grpc.CallOption) (QuerySvc_SubscribeToQueryStoreClient, error)
	UnsubscribeFromQueryStore(ctx context.Context, in *UnsubscribeFromQueryStoreRequest, opts ...grpc.CallOption) (*Empty, error)
}

type querySvcClient struct {
	cc grpc.ClientConnInterface
}

func NewQuerySvcClient(cc grpc.ClientConnInterface) QuerySvcClient {
	return &querySvcClient{cc}
}

func (c *querySvcClient) GetQuery(ctx context.Context, in *GetQueryRequest, opts ...grpc.CallOption) (*GetQueryResponse, error) {
	out := new(GetQueryResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.QuerySvc/GetQuery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *querySvcClient) GetQueries(ctx context.Context, in *GetQueriesRequest, opts ...grpc.CallOption) (*GetQueriesResponse, error) {
	out := new(GetQueriesResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.QuerySvc/GetQueries", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *querySvcClient) AddQuery(ctx context.Context, in *AddQueryRequest, opts ...grpc.CallOption) (*AddQueryResponse, error) {
	out := new(AddQueryResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.QuerySvc/AddQuery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *querySvcClient) UpdateQuery(ctx context.Context, in *UpdateQueryRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/m3uetcpb.QuerySvc/UpdateQuery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *querySvcClient) RemoveQuery(ctx context.Context, in *RemoveQueryRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/m3uetcpb.QuerySvc/RemoveQuery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *querySvcClient) QueryBy(ctx context.Context, in *QueryByRequest, opts ...grpc.CallOption) (*QueryByResponse, error) {
	out := new(QueryByResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.QuerySvc/QueryBy", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *querySvcClient) QueryInPlaylist(ctx context.Context, in *QueryInPlaylistRequest, opts ...grpc.CallOption) (*QueryInPlaylistResponse, error) {
	out := new(QueryInPlaylistResponse)
	err := c.cc.Invoke(ctx, "/m3uetcpb.QuerySvc/QueryInPlaylist", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *querySvcClient) QueryInQueue(ctx context.Context, in *QueryInQueueRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/m3uetcpb.QuerySvc/QueryInQueue", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *querySvcClient) SubscribeToQueryStore(ctx context.Context, in *Empty, opts ...grpc.CallOption) (QuerySvc_SubscribeToQueryStoreClient, error) {
	stream, err := c.cc.NewStream(ctx, &QuerySvc_ServiceDesc.Streams[0], "/m3uetcpb.QuerySvc/SubscribeToQueryStore", opts...)
	if err != nil {
		return nil, err
	}
	x := &querySvcSubscribeToQueryStoreClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type QuerySvc_SubscribeToQueryStoreClient interface {
	Recv() (*SubscribeToQueryStoreResponse, error)
	grpc.ClientStream
}

type querySvcSubscribeToQueryStoreClient struct {
	grpc.ClientStream
}

func (x *querySvcSubscribeToQueryStoreClient) Recv() (*SubscribeToQueryStoreResponse, error) {
	m := new(SubscribeToQueryStoreResponse)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *querySvcClient) UnsubscribeFromQueryStore(ctx context.Context, in *UnsubscribeFromQueryStoreRequest, opts ...grpc.CallOption) (*Empty, error) {
	out := new(Empty)
	err := c.cc.Invoke(ctx, "/m3uetcpb.QuerySvc/UnsubscribeFromQueryStore", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// QuerySvcServer is the server API for QuerySvc service.
// All implementations must embed UnimplementedQuerySvcServer
// for forward compatibility
type QuerySvcServer interface {
	GetQuery(context.Context, *GetQueryRequest) (*GetQueryResponse, error)
	GetQueries(context.Context, *GetQueriesRequest) (*GetQueriesResponse, error)
	AddQuery(context.Context, *AddQueryRequest) (*AddQueryResponse, error)
	UpdateQuery(context.Context, *UpdateQueryRequest) (*Empty, error)
	RemoveQuery(context.Context, *RemoveQueryRequest) (*Empty, error)
	QueryBy(context.Context, *QueryByRequest) (*QueryByResponse, error)
	QueryInPlaylist(context.Context, *QueryInPlaylistRequest) (*QueryInPlaylistResponse, error)
	QueryInQueue(context.Context, *QueryInQueueRequest) (*Empty, error)
	SubscribeToQueryStore(*Empty, QuerySvc_SubscribeToQueryStoreServer) error
	UnsubscribeFromQueryStore(context.Context, *UnsubscribeFromQueryStoreRequest) (*Empty, error)
	mustEmbedUnimplementedQuerySvcServer()
}

// UnimplementedQuerySvcServer must be embedded to have forward compatible implementations.
type UnimplementedQuerySvcServer struct {
}

func (UnimplementedQuerySvcServer) GetQuery(context.Context, *GetQueryRequest) (*GetQueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetQuery not implemented")
}
func (UnimplementedQuerySvcServer) GetQueries(context.Context, *GetQueriesRequest) (*GetQueriesResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetQueries not implemented")
}
func (UnimplementedQuerySvcServer) AddQuery(context.Context, *AddQueryRequest) (*AddQueryResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AddQuery not implemented")
}
func (UnimplementedQuerySvcServer) UpdateQuery(context.Context, *UpdateQueryRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateQuery not implemented")
}
func (UnimplementedQuerySvcServer) RemoveQuery(context.Context, *RemoveQueryRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method RemoveQuery not implemented")
}
func (UnimplementedQuerySvcServer) QueryBy(context.Context, *QueryByRequest) (*QueryByResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryBy not implemented")
}
func (UnimplementedQuerySvcServer) QueryInPlaylist(context.Context, *QueryInPlaylistRequest) (*QueryInPlaylistResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryInPlaylist not implemented")
}
func (UnimplementedQuerySvcServer) QueryInQueue(context.Context, *QueryInQueueRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method QueryInQueue not implemented")
}
func (UnimplementedQuerySvcServer) SubscribeToQueryStore(*Empty, QuerySvc_SubscribeToQueryStoreServer) error {
	return status.Errorf(codes.Unimplemented, "method SubscribeToQueryStore not implemented")
}
func (UnimplementedQuerySvcServer) UnsubscribeFromQueryStore(context.Context, *UnsubscribeFromQueryStoreRequest) (*Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UnsubscribeFromQueryStore not implemented")
}
func (UnimplementedQuerySvcServer) mustEmbedUnimplementedQuerySvcServer() {}

// UnsafeQuerySvcServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to QuerySvcServer will
// result in compilation errors.
type UnsafeQuerySvcServer interface {
	mustEmbedUnimplementedQuerySvcServer()
}

func RegisterQuerySvcServer(s grpc.ServiceRegistrar, srv QuerySvcServer) {
	s.RegisterService(&QuerySvc_ServiceDesc, srv)
}

func _QuerySvc_GetQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetQueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuerySvcServer).GetQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.QuerySvc/GetQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuerySvcServer).GetQuery(ctx, req.(*GetQueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuerySvc_GetQueries_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetQueriesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuerySvcServer).GetQueries(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.QuerySvc/GetQueries",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuerySvcServer).GetQueries(ctx, req.(*GetQueriesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuerySvc_AddQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AddQueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuerySvcServer).AddQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.QuerySvc/AddQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuerySvcServer).AddQuery(ctx, req.(*AddQueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuerySvc_UpdateQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateQueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuerySvcServer).UpdateQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.QuerySvc/UpdateQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuerySvcServer).UpdateQuery(ctx, req.(*UpdateQueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuerySvc_RemoveQuery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(RemoveQueryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuerySvcServer).RemoveQuery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.QuerySvc/RemoveQuery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuerySvcServer).RemoveQuery(ctx, req.(*RemoveQueryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuerySvc_QueryBy_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryByRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuerySvcServer).QueryBy(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.QuerySvc/QueryBy",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuerySvcServer).QueryBy(ctx, req.(*QueryByRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuerySvc_QueryInPlaylist_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryInPlaylistRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuerySvcServer).QueryInPlaylist(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.QuerySvc/QueryInPlaylist",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuerySvcServer).QueryInPlaylist(ctx, req.(*QueryInPlaylistRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuerySvc_QueryInQueue_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryInQueueRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuerySvcServer).QueryInQueue(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.QuerySvc/QueryInQueue",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuerySvcServer).QueryInQueue(ctx, req.(*QueryInQueueRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _QuerySvc_SubscribeToQueryStore_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(Empty)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(QuerySvcServer).SubscribeToQueryStore(m, &querySvcSubscribeToQueryStoreServer{stream})
}

type QuerySvc_SubscribeToQueryStoreServer interface {
	Send(*SubscribeToQueryStoreResponse) error
	grpc.ServerStream
}

type querySvcSubscribeToQueryStoreServer struct {
	grpc.ServerStream
}

func (x *querySvcSubscribeToQueryStoreServer) Send(m *SubscribeToQueryStoreResponse) error {
	return x.ServerStream.SendMsg(m)
}

func _QuerySvc_UnsubscribeFromQueryStore_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UnsubscribeFromQueryStoreRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(QuerySvcServer).UnsubscribeFromQueryStore(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/m3uetcpb.QuerySvc/UnsubscribeFromQueryStore",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(QuerySvcServer).UnsubscribeFromQueryStore(ctx, req.(*UnsubscribeFromQueryStoreRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// QuerySvc_ServiceDesc is the grpc.ServiceDesc for QuerySvc service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var QuerySvc_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "m3uetcpb.QuerySvc",
	HandlerType: (*QuerySvcServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetQuery",
			Handler:    _QuerySvc_GetQuery_Handler,
		},
		{
			MethodName: "GetQueries",
			Handler:    _QuerySvc_GetQueries_Handler,
		},
		{
			MethodName: "AddQuery",
			Handler:    _QuerySvc_AddQuery_Handler,
		},
		{
			MethodName: "UpdateQuery",
			Handler:    _QuerySvc_UpdateQuery_Handler,
		},
		{
			MethodName: "RemoveQuery",
			Handler:    _QuerySvc_RemoveQuery_Handler,
		},
		{
			MethodName: "QueryBy",
			Handler:    _QuerySvc_QueryBy_Handler,
		},
		{
			MethodName: "QueryInPlaylist",
			Handler:    _QuerySvc_QueryInPlaylist_Handler,
		},
		{
			MethodName: "QueryInQueue",
			Handler:    _QuerySvc_QueryInQueue_Handler,
		},
		{
			MethodName: "UnsubscribeFromQueryStore",
			Handler:    _QuerySvc_UnsubscribeFromQueryStore_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "SubscribeToQueryStore",
			Handler:       _QuerySvc_SubscribeToQueryStore_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "api/m3uetcpb/query.proto",
}
