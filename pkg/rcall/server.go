package rcall

import (
	"errors"
	"fmt"
	"go/token"
	"log"
	"net/http"
	"reflect"
)

type Server struct {
	services map[string]*service
}

func NewServer() *Server {
	return &Server{
		services: make(map[string]*service),
	}
}

type service struct {
	typ     reflect.Type
	recv    reflect.Value
	methods map[string]*serviceMethod
}
type serviceMethod struct {
	method    reflect.Method
	ArgType   reflect.Type
	ReplyType reflect.Type

	NeedReplyArg bool
	NeedContext  bool
}

func (svc *service) call(server *Server, req *Request, resp *Response, method *serviceMethod) {
	f := method.method.Func

	var argv = reflect.ValueOf(req.Argument)

	args := []reflect.Value{svc.recv, argv}
	if method.NeedReplyArg {
		replyv := reflect.New(method.ReplyType.Elem())
		switch method.ReplyType.Elem().Kind() {
		case reflect.Map:
			replyv.Elem().Set(reflect.MakeMap(method.ReplyType.Elem()))
		case reflect.Slice:
			replyv.Elem().Set(reflect.MakeSlice(method.ReplyType.Elem(), 0, 0))
		}
		args = append(args, replyv)
	}

	if method.NeedContext {
		args = append(args, reflect.ValueOf(req.Context))
	}

	ret := f.Call(args)

	if method.NeedReplyArg {
		resp.Reply = args[2].Interface()
	} else {
		resp.Reply = ret[0].Interface()
	}
	var err = ret[len(ret)-1]
	if !err.IsNil() {
		resp.Error = errorOfRemoteCall(err.Interface().(error))
	}
}

func errorOfRemoteCall(err error) *Error {
	if err == nil {
		return nil
	}
	return &Error{
		StatusCode: 500,
		Message:    err.Error(),
	}
}

func (svr *Server) ServeRequest(req *Request) (resp *Response) {
	resp = &Response{}

	svc, ok := svr.services[req.Coordinate.ToServicePath()]
	if !ok {
		s := fmt.Sprintf("rc.ServeRequest: service %q not found", req.Coordinate.ServiceName)
		log.Println(s)
		resp.Error = &Error{
			StatusCode: http.StatusBadRequest,
			Message:    s,
		}
		return
	}
	method, ok := svc.methods[req.MethodName]
	if !ok {
		s := fmt.Sprintf("rc.ServeRequest: service method not found %s.%s()", req.Coordinate.ServiceName, req.MethodName)
		log.Println(s)
		resp.Error = &Error{
			StatusCode: 401,
			Message:    s,
		}
	}

	svc.call(svr, req, resp, method)

	return
}

type ServiceRegisterConf struct {
	Target interface{}

	Coordinate ServiceCoordinate
}

func (svr *Server) MustRegister(conf ServiceRegisterConf) {
	if err := svr.Register(conf); err != nil {
		log.Fatal(err)
	}
}
func (svr *Server) Register(conf ServiceRegisterConf) error {
	return svr.register(conf)
}

func (svr *Server) register(conf ServiceRegisterConf) error {
	s := new(service)

	s.recv = reflect.ValueOf(conf.Target)

	typ := reflect.TypeOf(conf.Target)

	s.typ = typ

	coord := conf.Coordinate
	if coord.ServiceName == "" {
		coord.ServiceName = typ.Elem().Name()
	}

	if !token.IsExported(coord.ServiceName) {
		s := "rc.Register: type " + coord.ServiceName + " is not exported"
		log.Print(s)
		return errors.New(s)
	}

	log.Printf("rc.Register: register %s %s v%s", coord.Group, coord.ServiceName, coord.Version)
	s.methods = suitableMethods(typ, true)

	svr.services[coord.ToServicePath()] = s
	return nil
}

var typeOfError = reflect.TypeOf((*error)(nil)).Elem()

// suitableMethods supported method
// [ ] func ()(response,error)
// func (request)(response,error)
// func (request,*response)(error)
// func (request,context)(response,error)
// func (request,*response,context)(error)
func suitableMethods(typ reflect.Type, reportErr bool) map[string]*serviceMethod {
	methods := make(map[string]*serviceMethod)

	for m := 0; m < typ.NumMethod(); m++ {
		method := typ.Method(m)
		mtype := method.Type
		mname := method.Name
		// Method must be exported.
		if method.PkgPath != "" {
			continue
		}
		sm := &serviceMethod{
			method: method,
		}
		err := exractServiceMethod(sm)
		if err != nil {
			if reportErr {
				log.Println(err)
			}
			continue
		}

		log.Printf("rc.Register: add method %q: %s", mname, mtype.String())
		methods[mname] = sm
	}
	return methods
}

func exractServiceMethod(sm *serviceMethod) error {
	method := sm.method
	mtype := method.Type
	mname := method.Name

	numIn := mtype.NumIn()
	numOut := mtype.NumOut()
	if !(numIn > 1 && numOut >= 1) {
		goto invalidMethod
	}

	switch {
	// func (request)(response,error)
	case numIn == 2 && numOut == 2:
		sm.ArgType = mtype.In(1)
		sm.ReplyType = mtype.Out(0)
	// func (request,*response)(error)
	case numIn == 3 && numOut == 1:
		sm.ArgType = mtype.In(1)
		sm.ReplyType = mtype.In(2)
		sm.NeedReplyArg = true
	// func (request,context)(response,error)
	case numIn == 3 && numOut == 2:
		sm.ArgType = mtype.In(1)
		sm.ReplyType = mtype.In(2)
		sm.NeedContext = true

	// func (request,*response,context)(error)
	case numIn == 4 && numOut == 1:
		sm.ArgType = mtype.In(1)
		sm.ReplyType = mtype.Out(0)
		sm.NeedReplyArg = true
		sm.NeedContext = true
	default:
		goto invalidMethod
	}

	if sm.NeedReplyArg && sm.ReplyType.Kind() != reflect.Ptr {
		return fmt.Errorf("rc.Register: reply type of method %q is not a pointer: %q", mname, sm.ReplyType)
	}

	if returnType := mtype.Out(numOut - 1); returnType != typeOfError {
		return fmt.Errorf("rc.Register: return type of method %q is %q, must be error", mname, returnType)
	}

	return nil
invalidMethod:
	return fmt.Errorf("rc.Register: unsupported method %q: %s", mname, mtype.String())
}
