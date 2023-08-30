package middleware

import (
	"context"
	"fmt"
	"log/slog"
	"path"
	"time"

	"github.com/jwmwalrus/m3u-etcetera/internal/base"
	rtc "github.com/jwmwalrus/rtcycler"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type loggerMarker struct{}

var (
	loggerKey = &loggerMarker{}
)

type LoggerInterceptor struct {
	loggerType InterceptorType
	debug      bool
}

func NewLoggerInterceptor(lt InterceptorType, debug bool) *LoggerInterceptor {
	rtc.With(
		"type", lt,
		"debug", debug,
	).Trace("Creating Logger Interceptor")

	return &LoggerInterceptor{lt, debug}
}

func (li *LoggerInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{},
		info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {

		if path.Base(info.FullMethod) != "Off" {
			base.GetBusy(base.IdleStatusRequest)
			defer func() { base.GetFree(base.IdleStatusRequest) }()
		}

		newCtx := li.Before(ctx, info.FullMethod)
		res, err := handler(newCtx, req)
		li.After(newCtx, err)
		return res, err
	}
}

func (li *LoggerInterceptor) Stream() grpc.StreamServerInterceptor {
	return func(req interface{}, stream grpc.ServerStream,
		info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {

		base.GetBusy(base.IdleStatusRequest)
		defer func() { base.GetFree(base.IdleStatusRequest) }()

		newCtx := li.Before(stream.Context(), info.FullMethod)
		wrapped := newWrappedServerStream(newCtx, stream, li)
		wrapped.wrappedContext = newCtx
		err := handler(req, wrapped)
		li.After(newCtx, err)

		return err
	}
}

func (li *LoggerInterceptor) Before(ctx context.Context, fullMethod string) context.Context {

	start := time.Now()
	attrs := []slog.Attr{
		slog.Any("logger_type", li.loggerType),
		slog.Time("start", start),
		slog.String("grpc.service", path.Dir(fullMethod)[1:]),
		slog.String("grpc.method", path.Base(fullMethod)),
	}
	if d, ok := ctx.Deadline(); ok {
		attrs = append(attrs, slog.String("grpc.request.deadline", d.Format(time.RFC3339)))
	}

	newCtx := context.WithValue(ctx, loggerKey, attrs)
	return newCtx
}

func (li *LoggerInterceptor) After(ctx context.Context, err error) {
	finish := time.Now()
	bAttrs := ctx.Value(loggerKey).([]slog.Attr)
	s := status.Convert(err)
	code := s.Code()
	attrs := []slog.Attr{
		slog.Any("grpc.code", code),
	}

	var start time.Time
	for i := range bAttrs {
		switch bAttrs[i].Key {
		case "start":
			start = bAttrs[i].Value.Time()
			attrs = append(attrs, slog.String("grpc.start_time", start.Format(time.RFC3339)))
		default:
			attrs = append(attrs, bAttrs[i])
		}
	}
	diff := finish.Sub(start)
	attrs = append(attrs, slog.Float64("grpc.time_ms", float64(diff.Nanoseconds()/1000)/1000))
	if err != nil {
		attrs = append(attrs, slog.Any("error", err))
	}
	logw := slog.Default()
	msg := fmt.Sprintf(
		"Finished %v call with code %v",
		li.loggerType,
		code.String(),
	)

	switch code {
	case codes.OK,
		codes.Canceled,
		codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.Unauthenticated:
		if li.debug {
			logw.LogAttrs(ctx, slog.LevelDebug, msg, attrs...)
		} else {
			logw.LogAttrs(ctx, slog.LevelInfo, msg, attrs...)
		}
	case codes.DeadlineExceeded,
		codes.PermissionDenied,
		codes.ResourceExhausted,
		codes.FailedPrecondition,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unavailable:
		logw.LogAttrs(ctx, slog.LevelWarn, msg, attrs...)
	case codes.Unknown,
		codes.Unimplemented,
		codes.Internal,
		codes.DataLoss:
		logw.LogAttrs(ctx, slog.LevelError, msg, attrs...)
	default:
		logw.LogAttrs(ctx, slog.LevelError, msg, attrs...)
	}
}
