package rcall

import (
	"context"
	"fmt"
	"strconv"
)

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

type RemoteCallRequest struct {
	Coordinate ServiceCoordinate

	RequestID  string
	MethodName string

	Argument interface{}
	// GetArgument func(argv reflect.Value) (error)

	Context context.Context
	// Response *RemoteCallResponse
}

type RemoteCallResponse struct {
	RequestID string `json:"requestId"`
	Reply     interface{}
	// GetReply  func(ptr interface{}) error

	Error *RemoteCallError

	Context context.Context
}

type RemoteCallError struct {
	StatusCode int    `json:"statusCode"`
	ErrorCode  string `json:"code"`
	Message    string `json:"message"`
}

func (e RemoteCallError) Error() string {
	ec := e.ErrorCode
	if ec == "" {
		ec = strconv.Itoa(e.StatusCode)
	}
	return fmt.Sprintf("[%v/%v]: %s", ec, e.StatusCode, e.Message)
}
