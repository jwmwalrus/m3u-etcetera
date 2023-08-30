package middleware

import (
	"context"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type InterceptorType int

const (
	UnaryInterceptor InterceptorType = iota
	StreamInterceptor
)

func (it InterceptorType) String() string {
	return []string{
		"unary",
		"streaming",
	}[it]
}

type Interceptor interface {
	Unary() grpc.UnaryServerInterceptor
	Stream() grpc.StreamServerInterceptor
}

type WrappedInterceptor interface {
	Before(ctx context.Context, fullMethod string) context.Context
	After(ctx context.Context, err error)
}

// GetClientOpts returns the grpc dial options that any client should use.
func GetClientOpts() (opts []grpc.DialOption) {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return
}

// GetServerOpts returns the server middleware.
func GetServerOpts() (opts []grpc.ServerOption) {
	opts = append(
		opts,
		grpc.UnaryInterceptor(NewLoggerInterceptor(UnaryInterceptor, false).Unary()),
		grpc.StreamInterceptor(NewLoggerInterceptor(StreamInterceptor, false).Stream()),
	)

	return opts
}
