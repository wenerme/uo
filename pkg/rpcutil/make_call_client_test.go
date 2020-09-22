package rpcutil_test

import (
	"errors"
	"github.com/wenerme/uo/pkg/rpcutil"
	"net"
	"net/http"
	"net/rpc"
	"testing"
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
	rpc.Register(arith)
	rpc.HandleHTTP()
	l, e := net.Listen("tcp", ":1234")
	if e != nil {
		t.Fatal("listen error:", e)
	}
	go http.Serve(l, nil)

	// Client
	c := &ArithClient{}
	client, err := rpc.DialHTTP("tcp", "127.0.0.1:1234")
	if err != nil {
		panic(err)
	}

	if err := rpcutil.MakeCallClient(client.Call, "Arith", c); err != nil {
		t.Fatal(err)
	}

	if rel, err := c.Multiply(&Args{A: 10, B: 2}); err != nil {
		t.Fatal(err)
	} else {
		if rel != 20 {
			t.Fatal()
		}
	}
	if rel, err := c.Divide(&Args{A: 10, B: 2}); err != nil {
		t.Fatal(err)
	} else {
		if !(rel.Quo == 5 && rel.Rem == 0) {
			t.Fatal()
		}
	}
}
