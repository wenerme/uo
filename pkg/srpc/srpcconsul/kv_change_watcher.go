package srpcconsul

import (
	"sync"
	"time"

	kitlog "github.com/go-kit/kit/log"
	consulapi "github.com/hashicorp/consul/api"

	"github.com/wenerme/uo/pkg/kitutil"
)

type KeyValueChangeWatcher struct {
	Client       *consulapi.Client
	QueryOptions *consulapi.QueryOptions
	Key          string
	Wait         time.Duration

	Value func([]byte)
	Ready func()

	Logger kitlog.Logger

	c    chan struct{}
	lock sync.Mutex
}

func (r *KeyValueChangeWatcher) Watch() {
	r.lock.Lock()
	if r.c != nil {
		close(r.c)
	}
	r.c = make(chan struct{})
	r.lock.Unlock()

	wait := r.Wait
	if wait == 0 {
		wait = time.Second * 15
	}

	client := r.Client.KV()
	opts := r.QueryOptions
	key := r.Key
	log := kitutil.NewLogger(r.Logger)
	log = log.With("key", r.Key)

	lastIdx := uint64(0)

	cb := r.Value
	for {
		kv, qm, err := client.Get(key, opts)
		if err != nil {
			log.Warn("error", err, "KV Get failed")
		} else {
			changed := lastIdx == 0 || lastIdx != qm.LastIndex
			if changed {
				log.Info("event", "changed")
				if kv == nil {
					cb(nil)
				} else {
					cb(kv.Value)
				}
			}
			if lastIdx == 0 && r.Ready != nil {
				log.Info("event", "ready")
				r.Ready()
			}
			lastIdx = qm.LastIndex
		}

		select {
		case <-time.After(wait):
		case <-r.c:
			return
		}
	}
}
func (r *KeyValueChangeWatcher) Stop() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.c != nil {
		close(r.c)
		r.c = nil
	}
}
