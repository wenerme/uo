package querystring

import (
	"fmt"
	"net/url"
	"reflect"
	"time"

	"github.com/google/go-querystring/query"
)

var timeType = reflect.TypeOf(time.Time{})

// convert non struct to values
// for struct use https://github.com/google/go-querystring
func Values(v interface{}) (url.Values, error) {
	if v == nil {
		return nil, nil
	}
	switch tv := v.(type) {
	case map[string][]string:
		return tv, nil
	case map[string]string:
		return mapStringToSliceString(tv), nil
	}
	rv := reflect.ValueOf(v)
	if rv.Kind() == reflect.Struct {
		return query.Values(v)
	}
	if rv.Kind() == reflect.Map {
		m := make(url.Values)
		iter := rv.MapRange()
		for iter.Next() {
			k := iter.Key()
			sv := iter.Value()
			if sv.Kind() == reflect.Interface {
				sv = reflect.ValueOf(sv.Interface())
			}
			switch sv.Kind() {
			case
				reflect.Chan, reflect.Func, reflect.Map, reflect.Ptr, reflect.UnsafePointer,
				reflect.Interface, reflect.Slice:
				if sv.IsNil() {
					continue
				}
			}
			// no empty check
			if sv.Kind() == reflect.Slice || sv.Kind() == reflect.Array {
				for i := 0; i < sv.Len(); i++ {
					m.Add(fmt.Sprint(k.Interface()), valueString(sv.Index(i)))
				}
			} else {
				m.Set(fmt.Sprint(k.Interface()), valueString(sv))
			}
		}
		return m, nil
	}
	return nil, fmt.Errorf("httpmore.QueryValues: unsupported type %T", v)
}

func valueString(v reflect.Value) string {
	for v.Kind() == reflect.Ptr {
		if v.IsNil() {
			return ""
		}
		v = v.Elem()
	}

	if v.Type() == timeType {
		t := v.Interface().(time.Time)
		return t.Format(time.RFC3339)
	}

	return fmt.Sprint(v.Interface())
}

func mapStringToSliceString(a map[string]string) map[string][]string {
	if len(a) == 0 {
		return nil
	}
	m := make(map[string][]string)
	for k, v := range a {
		m[k] = []string{v}
	}
	return m
}
