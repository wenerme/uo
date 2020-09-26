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
	"fmt"
	stdlog "log"
	"net"
	"net/http"
	"strconv"

	"github.com/wenerme/uo/pkg/srpc/srpckit"

	"github.com/wenerme/uo/pkg/srpc/srpcconsul"

	consulapi "github.com/hashicorp/consul/api"

	kitlog "github.com/go-kit/kit/log"
	httptransport "github.com/go-kit/kit/transport/http"
	"github.com/gorilla/mux"
	"github.com/spf13/cobra"

	"github.com/wenerme/uo/cmd/pingsvc/pingapi"
	"github.com/wenerme/uo/pkg/kitutil"
	"github.com/wenerme/uo/pkg/srpc"
	"github.com/wenerme/uo/pkg/srpc/srpchttp"
)

var serverConf struct {
	HTTPAddress  string
	HTTPBind     string
	HTTPPort     int
	AdvertiseIP  string
	Consul       bool
	ConsulNodeID string
}

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start server",
	Long:  `PingService server`,
	Run: func(cmd *cobra.Command, args []string) {
		{
			noConsul, _ := cmd.Flags().GetBool("no-consul")
			serverConf.Consul = !noConsul

			var p string
			serverConf.HTTPBind, p, _ = net.SplitHostPort(serverConf.HTTPAddress)
			serverConf.HTTPPort, _ = strconv.Atoi(p)

			if serverConf.ConsulNodeID == "" {
				serverConf.ConsulNodeID = fmt.Sprintf("pingsvc-%s-%v", serverConf.AdvertiseIP, serverConf.HTTPPort)
			}
		}

		node := kitutil.MustMakeNodeContext(kitutil.NodeConf{
			Consul: serverConf.Consul,
		})
		_ = node.Logger.Log("phase", "start", "service", "ping")

		leader, err := node.ConsulClient.Status().Leader()
		_ = node.Logger.Log("msg", "consul leader", "leader", leader, "err", err)

		server := srpc.NewServer()

		pingService := &pingapi.PingService{}
		coordinate := srpc.GetCoordinate(pingService)
		server.MustRegister(srpc.ServiceRegisterConf{
			Target: pingService,
		})
		ep := srpckit.MakeServerEndpoint(server)
		ep = srpckit.InvokeLoggingMiddleware(kitlog.With(node.Logger, "server", "invoke"))(ep)
		serverHandler := httptransport.NewServer(ep, srpchttp.DecodeRequest, srpchttp.EncodeResponse)

		consulService := srpcconsul.SetServiceRegistration(coordinate, &consulapi.AgentServiceRegistration{
			ID:                serverConf.ConsulNodeID,
			Port:              serverConf.HTTPPort,
			EnableTagOverride: false,
			Checks: consulapi.AgentServiceChecks{
				&consulapi.AgentServiceCheck{
					Name:                           serverConf.ConsulNodeID + "-http-health",
					HTTP:                           fmt.Sprintf("http://127.0.0.1:%v/-/healthy", serverConf.HTTPPort),
					Interval:                       "15s",
					DeregisterCriticalServiceAfter: "30s",
				},
			},
		})

		sc, err := kitutil.MakeServiceEndpointContext(kitutil.ServiceEndpointConf{
			Node:          node,
			ConsulService: consulService,
		})
		if err != nil {
			panic(err)
		}
		sc.Registrar.Register()

		r := mux.NewRouter()
		r.Methods("POST").Path("/api/service/{group}/{service}/{version}/call/{method}").Handler(serverHandler)

		r.HandleFunc("/-/healthy", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(200)
			_, _ = rw.Write([]byte("OK"))
		})
		r.HandleFunc("/-/live", func(rw http.ResponseWriter, r *http.Request) {
			rw.WriteHeader(200)
			_, _ = rw.Write([]byte("OK"))
		})

		_ = node.Logger.Log("action", "start", "address", serverConf.HTTPAddress)
		stdlog.Fatal(http.ListenAndServe(serverConf.HTTPAddress, r))
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)

	serverCmd.Flags().StringVar(&serverConf.HTTPAddress, "http-address", "0.0.0.0:8123", "Listen host:port for HTTP endpoints")
	serverCmd.Flags().StringVar(&serverConf.AdvertiseIP, "advertise-ip", "127.0.0.1", "Sets the advertise address to use")
	serverCmd.Flags().StringVar(&serverConf.ConsulNodeID, "consul-node-id", "", "NodeID")
	serverCmd.Flags().Bool("no-consul", false, "Disable consul")
}
