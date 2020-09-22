package rpcutil_test

import (
	"errors"
	"net"
	"net/http"
	"net/rpc"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wenerme/uo/pkg/rpcutil"
)

type Args struct {
	A, B int
}

type Quotient struct {
	Quo, Rem int
}

type Arith int

func (t *Arith) Multiply(args *Args, reply *int) error {
	*reply = args.A * args.B
	return nil
}

func (t *Arith) Divide(args *Args, quo *Quotient) error {
	if args.B == 0 {
		return errors.New("divide by zero")
	}
	quo.Quo = args.A / args.B
	quo.Rem = args.A % args.B
	return nil
}

type ArithClient struct {
	Multiply func(args *Args) (int, error)
	Divide   func(args *Args) (Quotient, error)
}

func TestMakeCallClient(t *testing.T) {
	arith := new(Arith)
	assert.NoError(t, rpc.Register(arith))

	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	assert.NoError(t, e)

	go func() {
		assert.NoError(t, http.Serve(l, nil))
	}()

	// Client
	c := &ArithClient{}
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")

	assert.NoError(t, err)

	assert.NoError(t, rpcutil.MakeCallClient(client.Call, "Arith", c))

	{
		rel, err := c.Multiply(&Args{A: 10, B: 2})
		assert.NoError(t, err)
		assert.Equal(t, rel, 20)
	}

	{
		rel, err := c.Divide(&Args{A: 10, B: 2})
		assert.NoError(t, err)
		assert.Equal(t, rel, Quotient{Quo: 5, Rem: 0})
	}
}
