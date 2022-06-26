package middleware

import (
	"context"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"google.golang.org/grpc"
)

type wrappedServerStream struct {
	grpc.ServerStream
	wrappedContext context.Context
}

func streamServerInterceptor() grpc.StreamServerInterceptor {
	return func(req interface{}, stream grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		base.GetBusy(base.IdleStatusRequest)
		defer func() { base.GetFree(base.IdleStatusRequest) }()

		startTime := time.Now()

		newCtx := logBefore(
			stream.Context(),
			streamLogger,
			info.FullMethod,
			startTime,
		)
		wrapped := wrapServerStream(stream)
		wrapped.wrappedContext = newCtx
		err := handler(req, wrapped)
		logAfter(newCtx, err, time.Now(), false)

		return err
	}
}

func (w *wrappedServerStream) Context() context.Context {
	return w.wrappedContext
}

func (w *wrappedServerStream) SendMsg(m interface{}) error {
	err := w.ServerStream.SendMsg(m)
	logAfter(w.wrappedContext, err, time.Now(), true)
	return err
}

func wrapServerStream(stream grpc.ServerStream) *wrappedServerStream {
	if existing, ok := stream.(*wrappedServerStream); ok {
		return existing
	}
	return &wrappedServerStream{
		ServerStream:   stream,
		wrappedContext: stream.Context(),
	}
}
