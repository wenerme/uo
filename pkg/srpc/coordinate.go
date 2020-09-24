package srpc

import (
	"go/token"
	"reflect"
	"strings"
)

type hasServiceCoordinate interface {
	ServiceCoordinate() ServiceCoordinate
}

func GetCoordinate(v interface{}, override ServiceCoordinate) ServiceCoordinate {
	if sc, ok := v.(hasServiceCoordinate); ok {
		override = sc.ServiceCoordinate().Merge(override)
	}

	if override.ServiceName == "" {
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
			override.ServiceName = name
		}
	}
	return override.Normalize()
}

func (sc ServiceCoordinate) Merge(o ServiceCoordinate) ServiceCoordinate {
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
