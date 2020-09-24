package pingapi

import (
	"fmt"
	"time"

	"github.com/wenerme/uo/pkg/srpc"
)

type PingService struct {
}

func (s *PingService) Ping() string {
	return "PONG"
}
func (s *PingService) Hello(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}
func (s *PingService) Echo(v interface{}) (interface{}, error) {
	return v, nil
}
func (s *PingService) Now() time.Time {
	return time.Now()
}
func (s *PingService) ErrorOf(v string) (string, error) {
	return v, fmt.Errorf("ErrorOf: %s", v)
}
func (PingService) ServiceCoordinate() srpc.ServiceCoordinate {
	return srpc.ServiceCoordinate{
		ServiceName: "PingService",
		PackageName: "me.wener.ping",
	}
}

type PingServiceClient struct {
	Echo    func(v interface{}) (interface{}, error)
	ErrorOf func(msg string) (string, error)
}

func (PingServiceClient) ServiceCoordinate() srpc.ServiceCoordinate {
	return srpc.ServiceCoordinate{
		ServiceName: "PingService",
		PackageName: "me.wener.ping",
	}
}
