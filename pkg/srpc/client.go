package srpc

import (
	"context"
)

type Client struct {
	handler    InvokeFunc
	coordinate ServiceCoordinate
}

func NewClient(handler InvokeFunc, coordinate ServiceCoordinate) *Client {
	return &Client{
		handler:    handler,
		coordinate: coordinate,
	}
}

func (cli *Client) Call(ctx context.Context, methodName string, arg interface{}, reply interface{}) error {
	if ctx == nil {
		ctx = context.Background()
	}
	h := cli.handler
	r := h(&Request{
		Context:    ctx,
		MethodName: methodName,
		Coordinate: cli.coordinate,
		Argument:   arg,
	})

	if r.Error != nil {
		return r.Error
	}

	return r.GetReply(reply)
}

func (cli *Client) Close() error {
	return nil
}
