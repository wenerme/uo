package test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"testing"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"

	"github.com/wenerme/uo/pkg/srpc/httptrans"

	"github.com/davecgh/go-spew/spew"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/wenerme/uo/pkg/srpc"
)

type StringService struct {
}

func (*StringService) Uppercase(s string) (string, error) {
	return strings.ToUpper(s), nil
}

type StringServiceClient struct {
	Uppercase func(string) (string, error)
}

func TestSimpleCall(t *testing.T) {
	port := 8123
	{
		handler := makeTestServer()

		r := mux.NewRouter()
		r.Methods("POST").Path("/api/service/{group}/{service}/{version}/call/{method}").Handler(handler)

		httpServer := &http.Server{
			Handler: handler,
			Addr:    fmt.Sprintf(":%d", port),
		}
		go func() {
			_ = httpServer.ListenAndServe()
		}()
	}

	{
		options := []httptransport.ClientOption{
			httptransport.ClientBefore(httptrans.MakeRequestDumper(nil)),
			httptransport.ClientAfter(httptrans.MakeClientResponseDumper(nil)),
		}
		u, _ := url.Parse(fmt.Sprintf("http://localhost:%d", port))
		cli := httptransport.NewClient("POST", u, httptrans.EncodeRequest, httptrans.DecodeResponse, options...)
		ep := cli.Endpoint()

		client := &StringServiceClient{}

		assert.NoError(t, srpc.MakeRPCCallClient(func(ctx context.Context, request *srpc.Request) (response *srpc.Response, err error) {
			r, err := ep(ctx, request)
			if err != nil {
				log.Printf("Call failed %v", err)
				spew.Dump(r, err)
				return nil, err
			}
			return r.(*srpc.Response), err
		}, srpc.ServiceCoordinate{
			ServiceName: "StringService",
			PackageName: "com.example.test",
		}, client))

		reply, err := client.Uppercase("a")
		assert.NoError(t, err)
		assert.Equal(t, "A", reply)
	}
}

func makeTestServer() *httptransport.Server {
	svr := srpc.NewServer()
	svr.MustRegister(srpc.ServiceRegisterConf{
		Target: &StringService{},
		Coordinate: srpc.ServiceCoordinate{
			PackageName: "com.example.test",
		},
	})

	options := []httptransport.ServerOption{
		httptransport.ServerBefore(httptrans.MakeRequestDumper(nil)),
	}
	ep := srpc.MakeServerEndpoint(svr)
	logger := kitlog.NewLogfmtLogger(os.Stdout)
	ep = endpoint.Chain(srpc.InvokeLoggingMiddleware(logger))(ep)

	handler := httptransport.NewServer(
		ep,
		httptrans.DecodeRequest,
		httptrans.EncodeResponse,
		options...,
	)
	return handler
}
