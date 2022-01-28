package middleware

import (
	"google.golang.org/grpc"
)

// GetClientOpts returns the grpc dial options that any client should use
func GetClientOpts() (opts []grpc.DialOption) {
	opts = append(opts, grpc.WithInsecure())
	return
}

// GetServerOpts returns the server middleware
func GetServerOpts() (opts []grpc.ServerOption) {
	opts = append(
		opts,
		grpc.UnaryInterceptor(unaryServerInterceptor()),
		grpc.StreamInterceptor(streamServerInterceptor()),
	)

	return opts
}
