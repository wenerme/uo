package rcall

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/go-kit/kit/endpoint"
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

func LogInvokeMiddleware() endpoint.Middleware {
	return func(next endpoint.Endpoint) endpoint.Endpoint {
		return func(ctx context.Context, request interface{}) (interface{}, error) {
			if r, ok := request.(*Request); ok {
				st := time.Now()
				res, err := next(ctx, request)
				duration := time.Since(st)
				s := fmt.Sprintf("Invoke %s.%s %s/%s - %s ",
					r.Coordinate.ServiceTypeName(), r.MethodName, r.Coordinate.Group, r.Coordinate.Version, duration)
				if err == nil {
					if resp, ok := res.(*Response); ok && resp.Error != nil {
						err = resp.Error
					}
				}
				if err != nil {
					s += fmt.Sprintf("ERROR %v", err)
				} else {
					s += "OK"
				}
				log.Println(s)
				return res, err
			}

			return next(ctx, request)
		}
	}
}
