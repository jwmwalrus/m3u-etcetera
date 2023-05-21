package middleware

import (
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// GetClientOpts returns the grpc dial options that any client should use.
func GetClientOpts() (opts []grpc.DialOption) {
	opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	return
}

// GetServerOpts returns the server middleware.
func GetServerOpts() (opts []grpc.ServerOption) {
	opts = append(
		opts,
		grpc.UnaryInterceptor(unaryServerInterceptor()),
		grpc.StreamInterceptor(streamServerInterceptor()),
	)

	return opts
}
