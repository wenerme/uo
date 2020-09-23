package test

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"net/url"
	"strings"
	"testing"

	"github.com/wenerme/uo/pkg/rcall/httptrans"

	"github.com/davecgh/go-spew/spew"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/wenerme/uo/pkg/rcall"
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
		go assert.NoError(t, httpServer.ListenAndServe())
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

		assert.NoError(t, rcall.MakeRPCCallClient(func(ctx context.Context, request *rcall.RemoteCallRequest) (response *rcall.RemoteCallResponse, err error) {
			r, err := ep(ctx, request)
			if err != nil {
				log.Printf("Call failed %v", err)
				spew.Dump(r, err)
				return nil, err
			}
			return r.(*rcall.RemoteCallResponse), err
		}, rcall.ServiceCoordinate{
			ServiceName: "StringService",
			PackageName: "com.example.test",
		}, client))

		reply, err := client.Uppercase("a")
		assert.NoError(t, err)
		assert.Equal(t, "A", reply)
	}
}

func makeTestServer() *httptransport.Server {
	svr := rcall.NewServer()
	svr.MustRegister(rcall.ServiceRegisterConf{
		Target: &StringService{},
		Coordinate: rcall.ServiceCoordinate{
			PackageName: "com.example.test",
		},
	})

	options := []httptransport.ServerOption{
		httptransport.ServerBefore(httptrans.MakeRequestDumper(nil)),
	}
	handler := httptransport.NewServer(
		func(_ context.Context, request interface{}) (interface{}, error) {
			r := request.(*rcall.RemoteCallRequest)
			res := svr.ServeRequest(r)
			return res, nil
		},
		httptrans.DecodeRequest,
		httptrans.EncodeResponse,
		options...,
	)
	return handler
}
