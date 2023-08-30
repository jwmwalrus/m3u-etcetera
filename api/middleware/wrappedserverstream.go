package middleware

import (
	"context"

	"google.golang.org/grpc"
)

type wrappedServerStream struct {
	grpc.ServerStream
	wrappedContext context.Context
	interceptor    WrappedInterceptor
}

func newWrappedServerStream(ctx context.Context, stream grpc.ServerStream, si WrappedInterceptor) *wrappedServerStream {
	if existing, ok := stream.(*wrappedServerStream); ok {
		return existing
	}
	return &wrappedServerStream{
		ServerStream:   stream,
		wrappedContext: ctx,
		interceptor:    si,
	}
}

func (w *wrappedServerStream) Context() context.Context {
	return w.wrappedContext
}

func (w *wrappedServerStream) SendMsg(m interface{}) error {
	err := w.ServerStream.SendMsg(m)
	w.interceptor.After(w.wrappedContext, err)
	return err
}
