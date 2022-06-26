package middleware

import (
	"context"
	"path"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	"google.golang.org/grpc"
)

func unaryServerInterceptor() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		if path.Base(info.FullMethod) != "Off" {
			base.GetBusy(base.IdleStatusRequest)
			defer func() { base.GetFree(base.IdleStatusRequest) }()
		}

		startTime := time.Now()

		newCtx := logBefore(ctx, unaryLogger, info.FullMethod, startTime)
		res, err := handler(newCtx, req)
		logAfter(newCtx, err, time.Now(), false)
		return res, err
	}
}
