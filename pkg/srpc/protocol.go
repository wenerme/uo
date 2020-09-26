package srpc

import (
	"context"
	"fmt"

	"golang.org/x/mod/semver"
)

type HandlerFunc func(ctx context.Context, request *Request) (response *Response, err error)

var DefaultGroup = "default"
var DefaultVersion = "v1.0.0"

type ServiceCoordinate struct {
	Group       string
	Version     string // semver
	PackageName string
	ServiceName string
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
func (sc ServiceCoordinate) ServicePath() string {
	return fmt.Sprintf("%s/%s/%s", sc.Group, sc.ServiceTypeName(), sc.MajorVersion())
}
func (sc ServiceCoordinate) MajorVersion() string {
	if sc.Version == "" {
		return semver.Major(DefaultVersion)
	}
	return semver.Major(sc.Version)
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

	// Client use Argument to pass argument to server
	Argument interface{}
	// Server use GetArgument to get target typed replay
	GetArgument func(out interface{}) error

	Context context.Context
	// Response *RemoteCallResponse
}

type Response struct {
	RequestID string `json:"requestId"`
	// Server use Reply to return the response to client
	Reply interface{}
	// Client use GetReply to get target typed replay
	GetReply func(out interface{}) error

	Error *Error

	Context context.Context
}

type Error struct {
	StatusCode int    `json:"statusCode"`
	ErrorCode  int    `json:"code"`
	Message    string `json:"message"`
}

func (e Error) Error() string {
	ec := e.ErrorCode
	if ec == 0 {
		ec = e.StatusCode
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
