package kitutil

import (
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	consulapi "github.com/hashicorp/consul/api"
)

type ServiceEndpointConf struct {
	Node *NodeContext

	ConsulService *consulapi.AgentServiceRegistration
}

type ServiceEndpointContext struct {
	Registrar sd.Registrar
}

func MakeServiceEndpointContext(conf ServiceEndpointConf) (*ServiceEndpointContext, error) {
	node := conf.Node
	logger := node.Logger
	ctx := &ServiceEndpointContext{}

	if node.ConsulSdClient != nil && conf.ConsulService != nil {
		ctx.Registrar = consulsd.NewRegistrar(node.ConsulSdClient, conf.ConsulService, log.With(logger, "component", "consul-registrar"))
	}

	return ctx, nil
}
