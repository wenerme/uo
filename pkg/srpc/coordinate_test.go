package srpc_test

import (
	"testing"

	"github.com/magiconair/properties/assert"

	"github.com/wenerme/uo/pkg/srpc"
)

func TestNormalize(t *testing.T) {
	for _, test := range []struct {
		a srpc.ServiceCoordinate
		b srpc.ServiceCoordinate
	}{
		{
			srpc.ServiceCoordinate{},
			srpc.ServiceCoordinate{},
		},
		{
			srpc.ServiceCoordinate{ServiceName: "abc.d"},
			srpc.ServiceCoordinate{PackageName: "abc", ServiceName: "d"},
		},
	} {
		if test.b.Group == "" {
			test.b.Group = srpc.DefaultGroup
		}
		if test.b.Version == "" {
			test.b.Version = srpc.DefaultVersion
		}
		assert.Equal(t, test.a.Normalize(), test.b)
	}
}

func TestCoordinateMerge(t *testing.T) {
	a := srpc.ServiceCoordinate{ServiceName: "PingService"}
	b := a.WithOverride(srpc.ServiceCoordinate{Group: "auth"})
	assert.Equal(t, srpc.ServiceCoordinate{Group: "auth", ServiceName: "PingService"}, b)
	assert.Equal(t, srpc.ServiceCoordinate{ServiceName: "PingService"}, a)
}

func TestGetCoordinate(t *testing.T) {
	assert.Equal(t,
		srpc.GetCoordinate(&CoordinateTestService{}, srpc.ServiceCoordinate{PackageName: "me.wener"}),
		srpc.ServiceCoordinate{ServiceName: "CoordinateTestService", PackageName: "me.wener"}.Normalize(),
	)
	assert.Equal(t,
		srpc.GetCoordinate(&CoordinateTestServiceClient{}, srpc.ServiceCoordinate{Group: "auth", Version: "1.2.0"}),
		srpc.ServiceCoordinate{ServiceName: "CoordinateTestService", Group: "auth", Version: "1.2.0", PackageName: "me.wener.testing"},
	)
}

type CoordinateTestService struct {
}

func (CoordinateTestService) ServiceCoordinate() srpc.ServiceCoordinate {
	return srpc.ServiceCoordinate{
		PackageName: "me.wener.testing",
	}
}

type CoordinateTestServiceClient struct {
}

func (CoordinateTestServiceClient) ServiceCoordinate() srpc.ServiceCoordinate {
	return srpc.ServiceCoordinate{
		PackageName: "me.wener.testing",
	}
}
