package test

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/url"
	"os"
	"testing"
	"time"

	"github.com/go-kit/kit/endpoint"
	kitlog "github.com/go-kit/kit/log"

	"github.com/wenerme/uo/pkg/srpc/srpchttp"

	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/stretchr/testify/assert"

	"github.com/wenerme/uo/pkg/srpc"
)

func TestSimpleCall(t *testing.T) {
	port := 8123 + rand.Intn(10000)
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
			httptransport.ClientBefore(srpchttp.MakeRequestDumper(nil)),
			httptransport.ClientAfter(srpchttp.MakeClientResponseDumper(nil)),
		}
		u, _ := url.Parse(fmt.Sprintf("http://localhost:%d", port))
		cli := httptransport.NewClient("POST", u, srpchttp.EncodeRequest, srpchttp.DecodeResponse, options...)
		ep := cli.Endpoint()

		stringSvcClient := &StringServiceClient{}

		assert.NoError(t, srpc.MakeRPCCallClient(srpc.HandlerFuncOfEndpoint(ep), srpc.ServiceCoordinate{
			ServiceName: "StringService",
			PackageName: "com.example.test",
		}, stringSvcClient))

		StringServiceClientSpec(t, stringSvcClient)

		echoSvcClient := &EchoServiceClient{}
		assert.NoError(t, srpc.MakeRPCCallClient(srpc.HandlerFuncOfEndpoint(ep), srpc.ServiceCoordinate{
			ServiceName: "EchoService",
			PackageName: "com.example.test",
		}, echoSvcClient))

		EchoServiceClientSpec(t, echoSvcClient)
	}
}

func EchoServiceClientSpec(t *testing.T, echoSvcClient *EchoServiceClient) {
	for _, v := range []interface{}{
		"a",
		1,
		1.1,
		nil,
		map[string]interface{}{"Int": 1, "F": 1.1, "S": "ABC"},
		[]interface{}{"a", "b", "c", 1, 2, 3},
	} {
		replay, err := echoSvcClient.Echo(v)
		assert.NoError(t, err)
		assert.Equal(t, fmt.Sprintf("%#v", v), fmt.Sprintf("%#v", replay))
	}

	{
		v := time.Now()
		replay, err := echoSvcClient.EchoTime(v)
		assert.NoError(t, err)
		assert.True(t, v.Equal(replay))
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

	svr.MustRegister(srpc.ServiceRegisterConf{
		Target: &EchoService{},
		Coordinate: srpc.ServiceCoordinate{
			PackageName: "com.example.test",
		},
	})

	options := []httptransport.ServerOption{
		httptransport.ServerBefore(srpchttp.MakeRequestDumper(nil)),
	}
	ep := srpc.MakeServerEndpoint(svr)
	logger := kitlog.NewLogfmtLogger(os.Stdout)
	ep = endpoint.Chain(srpc.InvokeLoggingMiddleware(logger))(ep)

	handler := httptransport.NewServer(
		ep,
		srpchttp.DecodeRequest,
		srpchttp.EncodeResponse,
		options...,
	)
	return handler
}
