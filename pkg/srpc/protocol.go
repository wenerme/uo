package srpc

import (
	"context"
	"fmt"
)

type InvokeFunc func(request *Request) (response *Response)

var DefaultGroup = "default"
var DefaultVersion = "v1.0.0"

type ServiceCoordinate struct {
	Group       string
	Version     string // semver
	PackageName string
	ServiceName string
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
	RequestID string
	// Server use Reply to return the response to client
	Reply interface{}
	// Client use GetReply to get target typed replay
	GetReply func(out interface{}) error

	Error error

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
	switch t := err.(type) {
	case *Error:
		return t
	case Error:
		return &t
	}
	return &Error{
		StatusCode: 500,
		Message:    err.Error(),
	}
}

func ResponseOf(r *Request) *Response {
	return &Response{
		RequestID: r.RequestID,
		Context:   r.Context,
	}
}
