package pingapi

import (
	"context"
	"fmt"
	"time"
)

type PingService struct {
}

func (s *PingService) Ping() string {
	return "PONG"
}
func (s *PingService) Hello(name string) string {
	return fmt.Sprintf("Hello %s!", name)
}
func (s *PingService) Echo(ctx context.Context, v interface{}) (interface{}, error) {
	return v, nil
}
func (s *PingService) Now() time.Time {
	return time.Now()
}
func (s *PingService) ErrorOf(v string) (string, error) {
	return v, fmt.Errorf("ErrorOf: %s", v)
}

type PingServiceClient struct {
	ErrorOf func(msg string) (string, error)
}
