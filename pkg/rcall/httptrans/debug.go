package httptrans

import (
	"context"
	"log"
	"net/http"
	"net/http/httputil"

	httptransport "github.com/go-kit/kit/transport/http"
)

type DumperOptions struct {
	Disabled bool
	SkipBody bool
	Out      bool
	OnError  func(error)
}

func MakeRequestDumper(opt *DumperOptions) httptransport.RequestFunc {
	if opt == nil {
		opt = &DumperOptions{}
	}
	return func(ctx context.Context, r *http.Request) context.Context {
		var d []byte
		var err error
		if opt.Out {
			d, err = httputil.DumpRequestOut(r, !opt.SkipBody)
		} else {
			d, err = httputil.DumpRequest(r, !opt.SkipBody)
		}
		if err != nil {
			if opt.OnError != nil {
				opt.OnError(err)
			}
		} else {
			log.Println(string(d))
		}
		return ctx
	}
}

func MakeClientResponseDumper(opt *DumperOptions) httptransport.ClientResponseFunc {
	if opt == nil {
		opt = &DumperOptions{}
	}
	return func(ctx context.Context, r *http.Response) context.Context {
		d, err := httputil.DumpResponse(r, !opt.SkipBody)
		if err != nil {
			if opt.OnError != nil {
				opt.OnError(err)
			}
		} else {
			log.Println(string(d))
		}
		return ctx
	}
}
