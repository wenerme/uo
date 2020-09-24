package test

import (
	"time"

	"github.com/wenerme/uo/pkg/srpc"
)

type EchoService struct {
}

func (s *EchoService) Echo(v interface{}) (interface{}, error) {
	return v, nil
}
func (s *EchoService) EchoTime(v interface{}) (interface{}, error) {
	return v, nil
}

type EchoServiceClient struct {
	Echo     func(v interface{}) (interface{}, error)
	EchoTime func(v time.Time) (time.Time, error)
}

func (EchoService) ServiceCoordinate() srpc.ServiceCoordinate {
	return srpc.ServiceCoordinate{
		PackageName: "me.wener.test",
	}
}

func (EchoServiceClient) ServiceCoordinate() srpc.ServiceCoordinate {
	return srpc.ServiceCoordinate{
		PackageName: "me.wener.test",
	}
}
