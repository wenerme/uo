package kitutil

import (
	"time"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/log"
	"github.com/go-kit/kit/sd"
	consulsd "github.com/go-kit/kit/sd/consul"
	"github.com/go-kit/kit/sd/lb"
)

type BalanceStrategy string

const (
	RandomBalance     BalanceStrategy = "random"
	RoundRobinBalance BalanceStrategy = "rr"
)

type ClientEndpointConf struct {
	Node    *NodeContext
	Factory sd.Factory

	Instancer        bool
	InstancerService string
	InstancerTags    []string
	// Health only
	InstancerPassingOnly bool

	EndpointerOptions []sd.EndpointerOption

	BalancerStrategy   BalanceStrategy
	RandomBalancerSeed int64

	RetryMax      int
	RetryTimeout  time.Duration
	RetryCallback lb.Callback
}
type ClientEndpointContext struct {
	Endpoint   endpoint.Endpoint
	Endpointer sd.Endpointer
	Instancer  sd.Instancer
	Balancer   lb.Balancer

	ConsulInstancer *consulsd.Instancer
}

func MakeClientEndpointContext(conf ClientEndpointConf) (*ClientEndpointContext, error) {
	node := conf.Node
	logger := node.Logger

	ctx := &ClientEndpointContext{}

	// consul
	if node.ConsulSdClient != nil && conf.Instancer {
		ctx.ConsulInstancer = consulsd.NewInstancer(node.ConsulSdClient, log.With(logger, "component", "consul-instancer"), conf.InstancerService, conf.InstancerTags, conf.InstancerPassingOnly)
	}
	if ctx.ConsulInstancer != nil && ctx.Instancer == nil {
		ctx.Instancer = ctx.ConsulInstancer
	}

	// general
	if ctx.Instancer != nil && conf.Factory != nil {
		ctx.Endpointer = sd.NewEndpointer(ctx.Instancer, conf.Factory, log.With(logger, "component", "endpointer"), conf.EndpointerOptions...)
	}

	if ctx.Endpointer != nil {
		switch conf.BalancerStrategy {
		default:
			fallthrough
		case RoundRobinBalance:
			ctx.Balancer = lb.NewRoundRobin(ctx.Endpointer)
		case RandomBalance:
			ctx.Balancer = lb.NewRandom(ctx.Endpointer, conf.RandomBalancerSeed)
		}
	}

	if ctx.Balancer != nil && conf.RetryMax > 0 {
		if conf.RetryTimeout == 0 {
			conf.RetryTimeout = time.Minute
		}
		ctx.Endpoint = lb.Retry(conf.RetryMax, conf.RetryTimeout, ctx.Balancer)
	}
	if ctx.Balancer != nil && conf.RetryMax == 0 && conf.RetryCallback != nil {
		if conf.RetryTimeout == 0 {
			conf.RetryTimeout = time.Minute
		}
		ctx.Endpoint = lb.RetryWithCallback(conf.RetryTimeout, ctx.Balancer, conf.RetryCallback)
	}

	return ctx, nil
}
