package srpc

import (
	"context"
	"reflect"
)

var nilError = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())

// MakeRPCCallClient can make a struct as a rpc client, the method is defined as fields
func MakeRPCCallClient(handler InvokeFunc, coord ServiceCoordinate, v interface{}) error {
	val := reflect.ValueOf(v)
	typ := val.Type().Elem()

	coord = GetCoordinate(v).WithOverride(coord)

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if reflect.Func == f.Type.Kind() {
			ff := func(args []reflect.Value) (results []reflect.Value) {
				ev := nilError
				rv := reflect.New(f.Type.Out(0))

				req := &Request{
					Context:    context.Background(),
					Coordinate: coord,

					MethodName: f.Name,
					Argument:   args[0].Interface(),
				}

				resp := handler(req)

				if resp.Error != nil {
					return []reflect.Value{rv.Elem(), reflect.ValueOf(resp.Error)}
				}

				if err := resp.GetReply(rv.Interface()); err != nil {
					return []reflect.Value{rv.Elem(), reflect.ValueOf(err)}
				}

				results = []reflect.Value{rv.Elem(), ev}
				return
			}
			val.Elem().Field(i).Set(reflect.MakeFunc(f.Type, ff))
		}
	}
	return nil
}
