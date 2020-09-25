package srpcconsul

import (
	"fmt"
	"strings"

	"github.com/wenerme/uo/pkg/srpc/srpchttp"

	consulapi "github.com/hashicorp/consul/api"
	"golang.org/x/mod/semver"

	"github.com/wenerme/uo/pkg/srpc"
)

func GetServiceName(c srpc.ServiceCoordinate) string {
	return fmt.Sprintf("service-%s-%s", c.Group, strings.ReplaceAll(c.ServiceTypeName(), ".", "-"))
}

func GetServiceSelector(c srpc.ServiceCoordinate) (serviceName string, tags []string) {
	return GetServiceName(c), []string{
		"srpc",
		"version=" + semver.Major(c.Version),
	}
}

func SetServiceRegistration(c srpc.ServiceCoordinate, r *consulapi.AgentServiceRegistration) *consulapi.AgentServiceRegistration {
	if r.Name == "" {
		r.Name = GetServiceName(c)
	}

	r.Tags = append(r.Tags, "srpc",
		"group="+c.Group, "service="+c.ServiceName, "package="+c.PackageName,
		"name="+c.ServiceTypeName(),
		"version="+semver.Major(c.Version),
		// https://fabiolb.net/cfg/
		fmt.Sprintf("urlprefix-%s", srpchttp.ServicePrefix(c)),
	)

	if r.Meta == nil {
		r.Meta = make(map[string]string)
	}
	// meta version
	r.Meta["srpc"] = "1"
	r.Meta["group"] = c.Group
	r.Meta["name"] = c.ServiceTypeName()
	r.Meta["service"] = c.ServiceName
	r.Meta["package"] = c.PackageName
	r.Meta["version"] = c.Version

	return r
}
