package srpc_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/wenerme/uo/pkg/srpc"
	"github.com/wenerme/uo/pkg/srpc/test"
)

func TestServerErrorCase(t *testing.T) {
	server := srpc.NewServer()
	server.MustRegister(srpc.ServiceRegisterConf{
		Target: &test.StringService{},
	})

	{
		resp := server.ServeRequest(&srpc.Request{})
		assert.Equal(t, srpc.ErrCodeServiceNotFound, srpc.ErrorOf(resp.Error).ErrorCode)
	}
	{
		resp := server.ServeRequest(&srpc.Request{
			Coordinate: test.StringService{}.ServiceCoordinate(),
		})
		assert.Equal(t, srpc.ErrCodeMethodNotFound, srpc.ErrorOf(resp.Error).ErrorCode)
	}
}
