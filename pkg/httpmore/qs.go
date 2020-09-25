package httpmore

import (
	"fmt"
	"net/url"
	"reflect"
	"time"
)

var timeType = reflect.TypeOf(time.Time{})

// convert non struct to values
// for struct use https://github.com/google/go-querystring
func queryValues(v interface{}) (url.Values, error) {
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
	if rv.Kind() == reflect.Map {
		m := make(url.Values)
		iter := rv.MapRange()
		for iter.Next() {
			k := iter.Key()
			v := iter.Value()
			sv := reflect.ValueOf(v)

			if sv.Kind() == reflect.Slice || sv.Kind() == reflect.Array {
				for i := 0; i < sv.Len(); i++ {
					m.Add(fmt.Sprint(k), valueString(sv))
				}
			} else {
				m.Set(fmt.Sprint(k), valueString(sv))
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
