package httptrans

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"regexp"
	"strings"

	"github.com/gorilla/mux"
	jsoniter "github.com/json-iterator/go"

	"github.com/wenerme/uo/pkg/rcall"
)

const headerXRequestID = "X-Request-Id"

var regCall = regexp.MustCompile("/([^/]+)/([^/]+?)[.]([^/.]+)/([^/]+)/call/([^/]+)")

func DecodeRequest(ctx context.Context, r *http.Request) (interface{}, error) {
	vars := mux.Vars(r)

	if vars["service"] == "" {
		p := r.URL.Path
		prefix := "/api/service"
		if strings.HasPrefix(p, prefix) {
			m := regCall.FindStringSubmatch(p[len(prefix):])

			if m != nil {
				if vars == nil {
					vars = make(map[string]string)
				}
				vars["group"] = m[1]
				vars["package"] = m[2]
				vars["service"] = m[3]
				vars["version"] = m[4]
				vars["method"] = m[5]
			}
		}
	}

	req := &rcall.RemoteCallRequest{
		Coordinate: rcall.ServiceCoordinate{
			Group:       vars["group"],
			Version:     vars["version"],
			ServiceName: vars["service"],
			PackageName: vars["package"],
		},
		MethodName: vars["method"],
		Context:    ctx,
	}

	return req, jsoniter.NewDecoder(r.Body).Decode(&req.Argument)
}

func EncodeResponse(_ context.Context, rw http.ResponseWriter, res interface{}) error {
	r, ok := res.(*rcall.RemoteCallResponse)
	if !ok {
		return errors.New("httptrans.EncodeResponse: invalid response type")
	}

	rw.Header().Set("Content-Type", "application/json; charset=utf-8")

	if r.RequestID != "" {
		rw.Header().Set(headerXRequestID, r.RequestID)
	}
	if r.Error != nil {
		rw.WriteHeader(r.Error.StatusCode)
		return jsoniter.NewEncoder(rw).Encode(r.Error)
	}
	return jsoniter.NewEncoder(rw).Encode(r.Reply)
}

func EncodeRequest(ctx context.Context, req *http.Request, rcReq interface{}) error {
	r, ok := rcReq.(*rcall.RemoteCallRequest)
	if !ok {
		return errors.New("rc.EncodeRemoteCallRequest: invalid request type")
	}
	// fixme
	req.URL.Path = fmt.Sprintf("/api/service/%s/call/%s", r.Coordinate.ToServicePath(), r.MethodName)

	if r.RequestID != "" {
		req.Header.Set(headerXRequestID, r.RequestID)
	}
	req.Header.Set("Content-Type", "application/json; charset=utf-8")

	buf := &bytes.Buffer{}
	if err := jsoniter.NewEncoder(buf).Encode(r.Argument); err != nil {
		return err
	}

	req.Body = ioutil.NopCloser(buf)

	return nil
}

func DecodeResponse(ctx context.Context, resp *http.Response) (response interface{}, err error) {
	r := &rcall.RemoteCallResponse{
		Context: ctx,
	}
	r.RequestID = resp.Header.Get(headerXRequestID)

	if resp.StatusCode >= 400 {
		r.Error = &rcall.RemoteCallError{
			StatusCode: resp.StatusCode,
			Message:    resp.Status,
		}
		if err := jsoniter.NewDecoder(resp.Body).Decode(r.Error); err != nil {
			log.Printf("httptrans.DecodeResponse: encode error failed %s", err)
		}
		return nil, r.Error
	}

	if err = jsoniter.NewDecoder(resp.Body).Decode(&r.Reply); err != nil {
		return nil, err
	}

	return r, nil
}
