package rpcutil

import (
	"reflect"
)

var nilError = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())

// ClientCallHandler same as net/rpc client.Call
type ClientCallHandler func(serviceMethod string, args interface{}, reply interface{}) error

// MakeRPCCallClient can make a struct as a rpc client, the method is defined as fields
func MakeRPCCallClient(handler ClientCallHandler, serviceName string, v interface{}) error {
	val := reflect.ValueOf(v)
	typ := val.Type().Elem()

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if reflect.Func == f.Type.Kind() {
			ff := func(args []reflect.Value) (results []reflect.Value) {
				ev := nilError
				rv := reflect.New(f.Type.Out(0))

				err := handler(serviceName+"."+f.Name, args[0].Interface(), rv.Interface())
				if err != nil {
					ev = reflect.ValueOf(err)
				}

				results = []reflect.Value{rv.Elem(), ev}
				return
			}
			val.Elem().Field(i).Set(reflect.MakeFunc(f.Type, ff))
		}
	}
	return nil
}
