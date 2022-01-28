package middleware

import (
	"context"
	"fmt"
	"path"
	"time"

	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type loggerMarker struct{}

var (
	loggerKey = &loggerMarker{}
)

type loggerType int

const (
	unaryLogger loggerType = iota
	streamLogger
)

func (lt loggerType) String() string {
	return []string{
		"unary",
		"streaming",
	}[lt]
}

func logBefore(lt loggerType, ctx context.Context, fullMethod string, start time.Time) context.Context {
	fields := log.Fields{
		"loggerType":   lt,
		"start":        start,
		"grpc.service": path.Dir(fullMethod)[1:],
		"grpc.method":  path.Base(fullMethod),
	}
	if d, ok := ctx.Deadline(); ok {
		fields["grpc.request.deadline"] = d.Format(time.RFC3339)
	}

	newCtx := context.WithValue(ctx, loggerKey, fields)
	return newCtx
}

func logAfter(ctx context.Context, err error, finish time.Time, debug bool) {
	lf := ctx.Value(loggerKey).(log.Fields)
	s := status.Convert(err)
	code := s.Code()
	diff := finish.Sub(lf["start"].(time.Time))
	fields := log.Fields{
		"grpc.code":       code.String(),
		"grpc.start_time": lf["start"].(time.Time).Format(time.RFC3339),
		"grpc.service":    lf["grpc.service"],
		"grpc.method":     lf["grpc.method"],
		"grpc.time_ms":    float32(diff.Nanoseconds()/1000) / 1000,
	}
	if err != nil {
		fields[log.ErrorKey] = err
	}
	entry := log.WithContext(ctx).WithFields(fields)
	msg := fmt.Sprintf("Finished %v call with code %v", lf["loggerType"].(loggerType), code.String())

	switch code {
	case codes.OK,
		codes.Canceled,
		codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.Unauthenticated:
		if debug {
			entry.Debug(msg)
		} else {
			entry.Info(msg)
		}
	case codes.DeadlineExceeded,
		codes.PermissionDenied,
		codes.ResourceExhausted,
		codes.FailedPrecondition,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unavailable:
		entry.Warn(msg)
	case codes.Unknown,
		codes.Unimplemented,
		codes.Internal,
		codes.DataLoss:
		entry.Error(msg)
	default:
		entry.Error(msg)
	}
}

/*
func codeToLevel(code codes.Code) log.Level {
	switch code {
	case codes.OK,
		codes.Canceled,
		codes.InvalidArgument,
		codes.NotFound,
		codes.AlreadyExists,
		codes.Unauthenticated:
		return log.InfoLevel
	case codes.DeadlineExceeded,
		codes.PermissionDenied,
		codes.ResourceExhausted,
		codes.FailedPrecondition,
		codes.Aborted,
		codes.OutOfRange,
		codes.Unavailable:
		return log.WarnLevel
	case codes.Unknown,
		codes.Unimplemented,
		codes.Internal,
		codes.DataLoss:
		return log.ErrorLevel
	default:
		return log.ErrorLevel
	}
}
*/
