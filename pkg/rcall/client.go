package rcall

import (
	"context"
	"reflect"

	"github.com/mitchellh/mapstructure"
)

type RemoteCallHandlerFunc func(ctx context.Context, request *Request) (response *Response, err error)

var nilError = reflect.Zero(reflect.TypeOf((*error)(nil)).Elem())

// MakeRPCCallClient can make a struct as a rpc client, the method is defined as fields
func MakeRPCCallClient(handler RemoteCallHandlerFunc, coord ServiceCoordinate, v interface{}) error {
	val := reflect.ValueOf(v)
	typ := val.Type().Elem()

	coord = coord.Normalize()

	for i := 0; i < typ.NumField(); i++ {
		f := typ.Field(i)
		if reflect.Func == f.Type.Kind() {
			ff := func(args []reflect.Value) (results []reflect.Value) {
				ev := nilError
				rv := reflect.New(f.Type.Out(0))

				req := &Request{
					Coordinate: coord,

					MethodName: f.Name,
					Argument:   args[0].Interface(),
				}

				resp, err := handler(context.Background(), req)

				if err != nil {
					return []reflect.Value{rv.Elem(), reflect.ValueOf(err)}
				}

				if err := mapstructure.Decode(resp.Reply, rv.Interface()); err != nil {
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