/*
Copyright © 2020 wener <wenermail@gmail.com>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"context"
	"fmt"
	stdlog "log"

	"github.com/wenerme/uo/pkg/srpc/srpckit"

	"github.com/wenerme/uo/pkg/srpc/srpcconsul"

	"github.com/go-kit/kit/endpoint"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/spf13/cobra"

	"github.com/wenerme/uo/cmd/pingsvc/pingapi"
	"github.com/wenerme/uo/pkg/kitutil"
	"github.com/wenerme/uo/pkg/srpc"
	"github.com/wenerme/uo/pkg/srpc/srpchttp"
)

var callConf struct {
	HTTPAddress string
	Consul      bool

	Method string
	Arg    string
}

// callCmd represents the call command
var callCmd = &cobra.Command{
	Use:   "call",
	Short: "call ping service method",
	Long:  `PingService client`,
	Example: `
	# call Echo method
	go run . call -m Echo -a abc
`,
	Run: func(cmd *cobra.Command, args []string) {
		{
			noConsul, _ := cmd.Flags().GetBool("no-consul")
			serverConf.Consul = !noConsul
		}

		node := kitutil.MustMakeNodeContext(kitutil.NodeConf{
			Consul: serverConf.Consul,
		})
		_ = node.Logger.Log("phase", "start", "service", "ping")

		var ep endpoint.Endpoint
		factory := srpchttp.NewClientFactory(&srpchttp.ClientFactoryConf{
			Options: []httptransport.ClientOption{
				httptransport.ClientBefore(srpchttp.MakeRequestDumper(&srpchttp.DumperOptions{Out: true})),
				httptransport.ClientAfter(srpchttp.MakeClientResponseDumper(nil)),
			},
		})
		coordinate := srpc.GetCoordinate(pingapi.PingServiceClient{})
		sName, sTags := srpcconsul.GetServiceSelector(coordinate)
		if node.ConsulClient != nil {
			clientCtx, err := kitutil.MakeClientEndpointContext(kitutil.ClientEndpointConf{
				Node:             node,
				Factory:          factory,
				InstancerService: sName,
				InstancerTags:    sTags,
			})
			if err != nil {
				panic(err)
			}
			ep = clientCtx.Endpoint
		} else {
			edp, _, err := factory(callConf.HTTPAddress)
			if err != nil {
				panic(err)
			}
			ep = edp
		}

		ep = srpckit.InvokeLoggingMiddleware(node.Logger)(ep)
		client := &pingapi.PingServiceClient{}
		hfunc := srpckit.HandlerFuncOfEndpoint(ep)
		if err := srpc.MakeRPCCallClient(hfunc, srpc.ServiceCoordinate{}, client); err != nil {
			panic(err)
		}

		if callConf.Method == "" {
			stdlog.Fatal("Missing method name")
		}

		method := callConf.Method
		cli := srpc.NewClient(hfunc, client.ServiceCoordinate())
		var reply interface{}
		err := cli.Call(context.Background(), method, callConf.Arg, &reply)
		fmt.Printf("Method: %s\nArgument: %v\nError: %v\nReply: %v\n", method, callConf.Arg, err, reply)
	},
}

func init() {
	rootCmd.AddCommand(callCmd)

	callCmd.Flags().StringVar(&callConf.HTTPAddress, "http-address", "0.0.0.0:8123", "Listen host:port for HTTP endpoints")
	callCmd.Flags().StringVarP(&callConf.Method, "method", "m", "", "Listen host:port for HTTP endpoints")
	callCmd.Flags().StringVarP(&callConf.Arg, "arg", "a", "", "Listen host:port for HTTP endpoints")
	callCmd.Flags().Bool("no-consul", false, "Disable consul")
}
