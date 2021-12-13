package middleware

import "google.golang.org/grpc"

func GetServerOpts() (opts []grpc.ServerOption) {
	opts = append(opts, getLoggerOpts()...)
	return opts
}
