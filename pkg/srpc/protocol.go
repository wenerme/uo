package srpc

import (
	"context"
	"fmt"
	"strconv"
)

type HandlerFunc func(ctx context.Context, request *Request) (response *Response, err error)

var DefaultGroup = "default"
var DefaultVersion = "1.0.0"

type ServiceCoordinate struct {
	Group   string `json:"group"`
	Version string `json:"version"`

	PackageName string `json:"packageName"`
	ServiceName string `json:"serviceName"`
}

func (sc ServiceCoordinate) Normalize() ServiceCoordinate {
	g := sc.Group
	if g == "" {
		g = DefaultGroup
	}
	v := sc.Version
	if v == "" {
		v = DefaultVersion
	}
	return ServiceCoordinate{
		Group:       g,
		Version:     v,
		ServiceName: sc.ServiceName,
		PackageName: sc.PackageName,
	}
}
func (sc ServiceCoordinate) ToServicePath() string {
	c := sc.Normalize()
	return fmt.Sprintf("%s/%s/%s", c.Group, c.ServiceTypeName(), c.Version)
}
func (sc ServiceCoordinate) ServiceTypeName() string {
	if sc.PackageName != "" {
		return sc.PackageName + "." + sc.ServiceName
	}
	return sc.ServiceName
}

type Request struct {
	Coordinate ServiceCoordinate

	RequestID  string
	MethodName string

	Argument interface{}
	// GetArgument func(argv reflect.Value) (error)

	Context context.Context
	// Response *RemoteCallResponse
}

type Response struct {
	RequestID string `json:"requestId"`
	Reply     interface{}
	// GetReply  func(ptr interface{}) error

	Error *Error

	Context context.Context
}

type Error struct {
	StatusCode int    `json:"statusCode"`
	ErrorCode  string `json:"code"`
	Message    string `json:"message"`
}

func (e Error) Error() string {
	ec := e.ErrorCode
	if ec == "" {
		ec = strconv.Itoa(e.StatusCode)
	}
	return fmt.Sprintf("[%v/%v]: %s", ec, e.StatusCode, e.Message)
}

func ErrorOf(err error) *Error {
	if err == nil {
		return nil
	}
	if e, ok := err.(*Error); ok {
		return e
	}
	return &Error{
		StatusCode: 500,
		Message:    err.Error(),
	}
}
