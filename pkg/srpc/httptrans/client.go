package httptrans

import (
	"io"
	"net/url"
	"strings"

	"github.com/go-kit/kit/endpoint"
	"github.com/go-kit/kit/sd"
	httptransport "github.com/go-kit/kit/transport/http"
)

type ClientFactoryConf struct {
	Prefix  string
	Options []httptransport.ClientOption
}

func NewClientFactory(conf *ClientFactoryConf) sd.Factory {
	if conf == nil {
		conf = &ClientFactoryConf{}
	}
	prefix := conf.Prefix
	options := conf.Options
	return func(instance string) (endpoint.Endpoint, io.Closer, error) {
		if !strings.HasPrefix(instance, "http") {
			instance = "http://" + instance
		}
		u, err := url.Parse(instance + prefix)
		if err != nil {
			return nil, nil, err
		}
		cli := httptransport.NewClient("POST", u, EncodeRequest, DecodeResponse, options...)

		return cli.Endpoint(), nil, nil
	}
}
