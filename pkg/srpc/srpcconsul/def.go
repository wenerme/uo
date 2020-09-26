package srpcconsul

import (
	"fmt"
	"strings"

	"github.com/wenerme/uo/pkg/srpc/srpchttp"

	consulapi "github.com/hashicorp/consul/api"
	"github.com/wenerme/uo/pkg/srpc"
)

// Prefix of service kv store
func GetServiceKeyPrefix(c srpc.ServiceCoordinate) string {
	return fmt.Sprintf("services/" + c.ServicePath())
}

// Name of service
func GetServiceName(c srpc.ServiceCoordinate) string {
	return fmt.Sprintf("%s-%s-%s", c.Group, strings.ReplaceAll(c.ServiceTypeName(), ".", "-"), c.MajorVersion())
}

func GetServiceSelector(c srpc.ServiceCoordinate) (serviceName string, tags []string) {
	return GetServiceName(c), []string{
		"srpc",
		"version=" + c.MajorVersion(),
	}
}

func GetServiceTags(c srpc.ServiceCoordinate) []string {
	return []string{
		"srpc",
		"group=" + c.Group, "service=" + c.ServiceName, "package=" + c.PackageName,
		"name=" + c.ServiceTypeName(),
		"version=" + c.MajorVersion(),
		// https://fabiolb.net/cfg/
		fmt.Sprintf("urlprefix-%s", srpchttp.ServicePrefix(c)),
	}
}
func GetServiceMeta(c srpc.ServiceCoordinate) map[string]string {
	return map[string]string{
		"group":   c.Group,
		"name":    c.ServiceTypeName(),
		"service": c.ServiceName,
		"package": c.PackageName,
		"version": c.Version,
	}
}
func SetServiceRegistration(c srpc.ServiceCoordinate, r *consulapi.AgentServiceRegistration) *consulapi.AgentServiceRegistration {
	if r.Name == "" {
		r.Name = GetServiceName(c)
	}
	r.Tags = append(r.Tags, GetServiceTags(c)...)
	if r.Meta == nil {
		r.Meta = make(map[string]string)
	}
	// meta version
	r.Meta["srpc"] = "1"
	for k, v := range GetServiceMeta(c) {
		r.Meta[k] = v
	}
	return r
}
