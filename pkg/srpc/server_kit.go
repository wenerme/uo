package srpc

import (
	"context"
	"fmt"
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
)

func MakeServerEndpoint(svr *Server) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r, ok := request.(*Request)
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
			r, ok := request.(*Request)
			if !ok {
				return next(ctx, request)
			}

			st := time.Now()
			res, err := next(ctx, request)
			duration := time.Since(st)

			if err == nil {
				if resp, ok := res.(*Response); ok && resp.Error != nil {
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

func EndpointOfHandlerFunc(f HandlerFunc) endpoint.Endpoint {
	return func(ctx context.Context, request interface{}) (interface{}, error) {
		r, ok := request.(*Request)
		if !ok {
			return nil, fmt.Errorf("srpc.EndpointOfHandlerFunc: invalid request type %T", request)
		}

		return f(ctx, r)
	}
}
func HandlerFuncOfEndpoint(f endpoint.Endpoint) HandlerFunc {
	return func(ctx context.Context, request *Request) (*Response, error) {
		resp, err := f(ctx, request)
		if err != nil {
			return nil, err
		}

		r, ok := resp.(*Response)
		if !ok {
			return nil, fmt.Errorf("srpc.HandlerFuncOfEndpoint: invalid response type %T", resp)
		}
		return r, nil
	}
}
