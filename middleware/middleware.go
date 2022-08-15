package middleware

import (
	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/log"
	"github.com/rchauhan9/sportech/commons/go/logging"
)

func AddLogging(e endpoint.Endpoint, logger log.Logger) endpoint.Endpoint {
	return endpoint.Chain(
		logging.Middleware(logger),
	)(e)
}
