package srpc

import (
	"fmt"
	"go/token"
	"reflect"
	"strings"

	"golang.org/x/mod/semver"
)

type hasServiceCoordinate interface {
	ServiceCoordinate() ServiceCoordinate
}

func GetCoordinate(v interface{}) ServiceCoordinate {
	c := ServiceCoordinate{}
	if sc, ok := v.(hasServiceCoordinate); ok {
		c = sc.ServiceCoordinate()
	}

	if c.ServiceName == "" && v != nil {
		rt := reflect.TypeOf(v)
		if rt.Kind() == reflect.Ptr {
			rt = rt.Elem()
		}
		name := rt.Name()
		if token.IsExported(name) {
			// StringServiceClient -> StringService
			if strings.HasSuffix(name, "Client") {
				name = name[:len(name)-len("Client")]
			}
			c.ServiceName = name
		}
	}
	return c.Normalize()
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
	p := sc.PackageName
	s := sc.ServiceName
	if p == "" && strings.Contains(s, ".") {
		i := strings.LastIndex(s, ".")
		p = s[:i]
		s = s[i+1:]
	}
	return ServiceCoordinate{
		Group:       g,
		Version:     v,
		PackageName: p,
		ServiceName: s,
	}
}

func (sc ServiceCoordinate) ServicePath() string {
	g := sc.Group
	if g == "" {
		g = DefaultGroup
	}
	return fmt.Sprintf("%s/%s/%s", g, sc.ServiceTypeName(), sc.MajorVersion())
}
func (sc ServiceCoordinate) MajorVersion() string {
	if sc.Version == "" {
		return semver.Major(DefaultVersion)
	}
	return semver.Major(sc.Version)
}
func (sc ServiceCoordinate) ServiceTypeName() string {
	if sc.PackageName != "" {
		return sc.PackageName + "." + sc.ServiceName
	}
	return sc.ServiceName
}

func (sc ServiceCoordinate) WithOverride(o ServiceCoordinate) ServiceCoordinate {
	if o.ServiceName != "" {
		sc.ServiceName = o.ServiceName
	}
	if o.PackageName != "" {
		sc.PackageName = o.PackageName
	}
	if o.Version != "" && o.Version != DefaultVersion {
		sc.Version = o.Version
	}
	if o.Group != "" && o.Group != DefaultGroup {
		sc.Group = o.Group
	}
	return sc
}

func (sc ServiceCoordinate) IsValid() bool {
	return sc.ServiceName != ""
}
