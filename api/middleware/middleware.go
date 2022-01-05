package middleware

import (
	"google.golang.org/grpc"
)

// GetServerOpts returns the server middleware
func GetServerOpts() (opts []grpc.ServerOption) {
	// opts = append(opts, getLoggerOpts()...)

	opts = append(
		opts,
		grpc.UnaryInterceptor(unaryServerInterceptor()),
		grpc.StreamInterceptor(streamServerInterceptor()),
	)

	return opts
}
