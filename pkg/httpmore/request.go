package httpmore

import (
	"bytes"
	"context"
	"fmt"
	"io"
	"io/ioutil"
	stdlog "log"
	"net/http"
	"net/url"
	"strings"

	"github.com/wenerme/uo/querystring"

	jsoniter "github.com/json-iterator/go"
)

type RequestInit struct {
	Method   string
	BaseURL  string
	URL      string
	Query    interface{}
	Body     io.Reader
	JSONBody interface{}
	Header   http.Header
	Context  context.Context
	Options  Values // Extra options for customized process - non string option use Context
}

func (r RequestInit) WithOverride(o RequestInit) RequestInit {
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

	switch {
	case r.Query == nil:
		r.Query = o.Query
	case o.Query == nil:
		// keep
	default:
		if a, ae := querystring.Values(r.Query); ae == nil {
			if b, be := querystring.Values(o.Query); be == nil {
				r.Query = mergeMapSliceString(a, b)
			} else {
				stdlog.Printf("httmore.RequestInit.Merge: convert query failed %v", be)
			}
		} else {
			stdlog.Printf("httmore.RequestInit.Merge: convert query failed %v", ae)
			r.Query = o.Query
		}
	}

	r.Header = mergeMapSliceString(r.Header, o.Header)
	if r.Options == nil {
		r.Options = Values{}
	}
	r.Options = r.Options.Clone().WithMerge(o.Options)
	return r
}

func (r RequestInit) NewRequest() (*http.Request, error) {
	ctx := r.Context
	if ctx == nil {
		ctx = context.Background()
	}
	u, err := r.GetURL()
	if err != nil {
		return nil, err
	}
	req, err := http.NewRequestWithContext(ctx, r.Method, u, nil)
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
func (r *RequestInit) GetURL() (string, error) {
	u := r.URL
	if strings.HasPrefix(u, "/") {
		u = r.BaseURL + u
	}
	v, err := querystring.Values(r.Query)
	if err != nil {
		return "", err
	}
	if len(v) > 0 {
		if parsed, err := url.Parse(u); err == nil {
			parsed.RawQuery = (url.Values)(mergeMapSliceString(v, parsed.Query())).Encode()
			u = parsed.String()
		} else {
			return "", fmt.Errorf("httpmore.GetURL: parse url failed %v", err)
		}
	}
	return u, nil
}
