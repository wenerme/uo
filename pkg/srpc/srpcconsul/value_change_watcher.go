package srpcconsul

import (
	"bytes"
	"reflect"
	"sync"
	"time"

	kitlog "github.com/go-kit/kit/log"
	jsoniter "github.com/json-iterator/go"

	"github.com/wenerme/uo/pkg/kitutil"
)

type StateChangeWatcher struct {
	Value interface{}
	Wait  time.Duration

	Change func([]byte)
	Ready  func()

	Logger kitlog.Logger

	c    chan struct{}
	lock sync.Mutex
}

func (r *StateChangeWatcher) Watch() {
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

	log := kitutil.NewLogger(r.Logger)
	log = log.With("type", reflect.TypeOf(r.Value).String())

	var lastData []byte
	for {
		val := r.Value
		data, err := jsoniter.Marshal(val)
		if err != nil {
			log.Warn("error", err, "msg", "failed to Marshal")
		} else {
			if !bytes.Equal(data, lastData) {
				log.Info("event", "changed")
				r.Change(data)
			}
			if lastData == nil && r.Ready != nil {
				log.Info("event", "ready")
				r.Ready()
			}
			lastData = data
		}

		select {
		case <-time.After(wait):
		case <-r.c:
			return
		}
	}
}
func (r *StateChangeWatcher) Stop() {
	r.lock.Lock()
	defer r.lock.Unlock()
	if r.c != nil {
		close(r.c)
		r.c = nil
	}
}
