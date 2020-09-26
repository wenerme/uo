package srpcconsul

import (
	"errors"
	"fmt"
	"sync"
	"time"

	jsoniter "github.com/json-iterator/go"

	"github.com/wenerme/uo/pkg/kitutil"

	kitlog "github.com/go-kit/kit/log"
	consulapi "github.com/hashicorp/consul/api"
	"github.com/oklog/run"

	"github.com/wenerme/uo/pkg/srpc"
)

type Configurer interface {
	Run() error // blocked run
	Ready() <-chan struct{}
	IsReady() bool
	Stop()
}
type configurer struct {
	g      run.Group
	readwg sync.WaitGroup
	readyc chan struct{}
	ready  bool
	c      chan struct{}
}

func (c *configurer) Run() error {
	return c.g.Run()
}
func (c *configurer) Ready() <-chan struct{} {
	return c.readyc
}
func (c *configurer) IsReady() bool {
	return c.ready
}
func (c *configurer) Stop() {
	close(c.c)
}

type ServiceConfigurerOptions struct {
	Client            *consulapi.Client
	Service           interface{}
	ServiceCoordinate srpc.ServiceCoordinate // override
	Logger            kitlog.Logger

	Conf  interface{}
	State interface{}
}

func addStateConfigurer(cfg *configurer, opts ServiceConfigurerOptions, coord srpc.ServiceCoordinate, log kitutil.Logger) {
	cfg.readwg.Add(1)
	key := fmt.Sprintf("%s/state", GetServiceKeyPrefix(coord))
	w := &StateChangeWatcher{
		Logger: log.With("comp", "StateConfigurer"),
		Ready: func() {
			cfg.readwg.Done()
		},
		Value: opts.State,
		Change: func(bytes []byte) {
			if bytes != nil {
				_, err := opts.Client.KV().Put(&consulapi.KVPair{
					Key:   key,
					Value: bytes,
				}, nil)
				if err != nil {
					log.Error("error", err, "action", "PutState")
				}
			}
		},
	}
	cfg.g.Add(func() error {
		// read initial
		for {
			kv, _, err := opts.Client.KV().Get(key, nil)
			if err == nil {
				if kv != nil {
					err := jsoniter.Unmarshal(kv.Value, opts.State)
					if err != nil {
						log.Error("error", err, "action", "UnmarshalInitialState")
					}
				}
				break
			}
			log.Error("error", err, "action", "GetInitialState")
			time.Sleep(time.Second * 5)
		}

		w.Watch()
		return nil
	}, func(err error) {
		w.Stop()
	})
}
func addConfConfigurer(cfg *configurer, opts ServiceConfigurerOptions, coord srpc.ServiceCoordinate, log kitutil.Logger) {
	cfg.readwg.Add(1)
	key := fmt.Sprintf("%s/conf", GetServiceKeyPrefix(coord))
	w := &KeyValueChangeWatcher{
		Client: opts.Client,
		Logger: log.With("comp", "ConfConfigurer"),
		Key:    key,
		Ready: func() {
			cfg.readwg.Done()
		},
		Value: func(data []byte) {
			if data != nil {
				err := jsoniter.Unmarshal(data, opts.Conf)
				if err != nil {
					log.Warn("error", err, "action", "UnmarshalConf")
				}
			}
		},
	}
	cfg.g.Add(func() error {
		// put initial
		for {
			kv, _, err := opts.Client.KV().Get(key, nil)
			if err == nil {
				if kv == nil {
					b, err := jsoniter.Marshal(opts.Conf)
					if err != nil {
						log.Error("error", err, "action", "MarshalInitialConf")
						break
					}

					if _, err := opts.Client.KV().Put(&consulapi.KVPair{
						Key:   key,
						Value: b,
					}, nil); err != nil {
						log.Error("error", err, "action", "PutInitialConf")
					} else {
						log.Info("action", "PutInitialConf")
					}
				}
				break
			}
			log.Error("error", err, "action", "PutInitialConf")
			time.Sleep(time.Second * 5)
		}
		w.Watch()
		return nil
	}, func(err error) {
		w.Stop()
	})
}
func NewServiceConfigurer(opts ServiceConfigurerOptions) (Configurer, error) {
	cfg := &configurer{}
	coord := srpc.GetCoordinate(opts.Service).WithOverride(opts.ServiceCoordinate)
	if !coord.IsValid() {
		return nil, errors.New("invalid service coordinate")
	}
	log := kitutil.NewLogger(opts.Logger)
	log = log.With("service", coord.ServiceTypeName())
	n := 0
	if opts.Conf != nil {
		n++
		addConfConfigurer(cfg, opts, coord, log)
	}
	if opts.State != nil {
		n++
		addStateConfigurer(cfg, opts, coord, log)
	}
	if n == 0 {
		return nil, errors.New("nothing to config")
	}

	cfg.g.Add(func() error {
		<-cfg.c
		return errors.New("closed")
	}, func(err error) {
		if err != nil && err.Error() != "closed" {
			close(cfg.c)
		}
	})

	cfg.c = make(chan struct{})
	cfg.readyc = make(chan struct{})
	go func() {
		cfg.readwg.Wait()
		cfg.ready = true
		close(cfg.readyc)
	}()

	return cfg, nil
}
