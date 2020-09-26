package srpchttp

import (
	"fmt"

	"github.com/wenerme/uo/pkg/srpc"
)

func ServicePrefix(c srpc.ServiceCoordinate) string {
	return fmt.Sprintf("%s/%s", DefaultPrefix, c.ServicePath())
}
