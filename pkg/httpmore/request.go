package httpmore

import (
	"bytes"
	"context"
	"io"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"net/url"
	"strings"

	jsoniter "github.com/json-iterator/go"
)

type RequestInit struct {
	Method   string
	BaseURL  string
	URL      string
	Query    url.Values
	Body     io.Reader
	JSONBody interface{}
	Header   http.Header
	Context  context.Context
}

func (r RequestInit) Merge(o RequestInit) RequestInit {
	if o.Method != "" {
		r.Method = o.Method
	}
	if o.BaseURL != "" {
		r.BaseURL = o.BaseURL
	}
	if o.URL != "" {
		r.URL = o.URL
	}

	if o.Body != nil {
		r.Body = o.Body
	}
	if o.Context != nil {
		r.Context = o.Context
	}

	r.Query = mergeMapSliceString(r.Query, o.Query)
	r.Header = mergeMapSliceString(r.Header, o.Header)

	return r
}

func (r RequestInit) NewRequest() (*http.Request, error) {
	ctx := r.Context
	if ctx == nil {
		ctx = context.Background()
	}
	req, err := http.NewRequestWithContext(ctx, r.Method, r.GetURL(), nil)
	if err != nil {
		return nil, err
	}
	if len(r.Header) > 0 {
		req.Header = r.Header
	}
	if r.HasBody() {
		b, err := r.GetBody()
		if err != nil {
			return nil, err
		}
		req.Body = b
		req.GetBody = r.GetBody
	}

	return req, nil
}
func (r *RequestInit) HasBody() bool {
	return r.Body != nil || r.JSONBody != nil
}
func (r *RequestInit) GetBody() (io.ReadCloser, error) {
	b := r.Body
	if b == nil && r.JSONBody != nil {
		buf := &bytes.Buffer{}
		err := jsoniter.NewEncoder(buf).Encode(r)
		if err != nil {
			return nil, err
		}
		b = buf
	}
	if rc, ok := b.(io.ReadCloser); ok {
		return rc, nil
	}
	return ioutil.NopCloser(b), nil
}
func (r *RequestInit) GetURL() string {
	u := r.URL
	if strings.HasPrefix(u, "/") {
		u = r.BaseURL + u
	}
	if len(r.Query) > 0 {
		if parsed, err := url.Parse(u); err == nil {
			parsed.RawQuery = (url.Values)(mergeMapSliceString(r.Query, parsed.Query())).Encode()
			u = parsed.String()
		} else {
			stdlog.Printf("httpmore.GetURL: parse url failed %v", err)
		}
	}
	return u
}
