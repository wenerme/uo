package kitutil

import (
	stdlog "log"
	"os"

	"github.com/go-kit/kit/log"
	consulsd "github.com/go-kit/kit/sd/consul"
	consulapi "github.com/hashicorp/consul/api"
)

type NodeConf struct {
	Logger     log.Logger
	Consul     bool
	ConsulConf *consulapi.Config
}
type NodeContext struct {
	Logger       log.Logger
	ConsulClient *consulapi.Client

	ConsulSdClient consulsd.Client
}

func MustMakeNodeContext(conf NodeConf) *NodeContext {
	v, err := MakeNodeContext(conf)
	if err != nil {
		stdlog.Fatal(err)
	}
	return v
}
func MakeNodeContext(conf NodeConf) (*NodeContext, error) {
	logger := conf.Logger
	if logger == nil {
		logger = log.NewLogfmtLogger(os.Stderr)
		logger = log.With(logger, "ts", log.DefaultTimestampUTC)
		logger = log.With(logger, "caller", log.DefaultCaller)
	}

	ctx := &NodeContext{
		Logger: logger,
	}

	if conf.Consul {
		cc := conf.ConsulConf
		if cc == nil {
			cc = consulapi.DefaultConfig()
		}

		cli, err := consulapi.NewClient(cc)
		if err != nil {
			return nil, err
		}
		ctx.ConsulClient = cli
	}

	if ctx.ConsulClient != nil {
		ctx.ConsulSdClient = consulsd.NewClient(ctx.ConsulClient)
	}

	return ctx, nil
}
