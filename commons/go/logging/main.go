package logging

import (
	"context"
	"reflect"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/go-kit/log/level"
)

func Middleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (_ interface{}, err error) {
			defer func(startTime time.Time) {
				responseTime := time.Since(startTime)
				level.Info(logger).Log("request", reflect.TypeOf(request), "time", responseTime, "err", err)
			}(time.Now())
			return next(ctx, request)
		}
	}
}
