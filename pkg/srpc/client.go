package srpc

import (
	"context"
)

type Client struct {
	handler    HandlerFunc
	coordinate ServiceCoordinate
}

func NewClient(handler HandlerFunc, coordinate ServiceCoordinate) *Client {
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
	r, err := h(ctx, &Request{
		Context:    ctx,
		MethodName: methodName,
		Coordinate: cli.coordinate,
		Argument:   arg,
	})

	if err != nil {
		return err
	}
	return r.GetReply(reply)
}

func (cli *Client) Close() error {
	return nil
}
