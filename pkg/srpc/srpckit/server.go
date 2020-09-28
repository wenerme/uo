package srpckit

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"

	"github.com/wenerme/uo/pkg/srpc"
)

func MakeServerEndpoint(svr *srpc.Server) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r, ok := request.(*srpc.Request)
		if !ok {
			return nil, fmt.Errorf("rcall.MakeServerEndpoint: invalid request type %T", request)
		}
		if r.Context == nil {
			r.Context = ctx
		}
		res := svr.ServeRequest(r)
		return res, nil
	}
}

func InvokeLoggingMiddleware(logger log.Logger) endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			r, ok := request.(*srpc.Request)
			if !ok {
				return next(ctx, request)
			}

			st := time.Now()
			res, err := next(ctx, request)
			duration := time.Since(st)

			if err == nil {
				if resp, ok := res.(*srpc.Response); ok && resp.Error != nil {
					err = resp.Error
				}
			}

			_ = logger.Log("service", r.Coordinate.ServiceTypeName(), "method", r.MethodName,
				"group", r.Coordinate.Group, "version", r.Coordinate.Version,
				"time", duration, "err", err,
			)

			return res, err
		}
	}
}

func EndpointOfHandlerFunc(f srpc.InvokeFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r, ok := request.(*srpc.Request)
		if !ok {
			return nil, fmt.Errorf("srpc.EndpointOfHandlerFunc: invalid request type %T", request)
		}
		if r.Context == nil {
			r.Context = ctx
		}
		return f(r), nil
	}
}
func HandlerFuncOfEndpoint(f endpoint.Endpoint) srpc.InvokeFunc {
	return func(request *srpc.Request) *srpc.Response {
		resp, err := f(request.Context, request)
		if err != nil {
			r := srpc.ResponseOf(request)
			r.Error = err
			return r
		}

		r, ok := resp.(*srpc.Response)
		if !ok {
			r := srpc.ResponseOf(request)
			r.Error = fmt.Errorf("srpc.HandlerFuncOfEndpoint: invalid response type %T", resp)
			return r
		}
		return r
	}
}
