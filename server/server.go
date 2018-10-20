package server

import (
	"fmt"
	"net/http"
	"reflect"
	"strings"

	"github.com/dynamicgo/slf4go"

	"github.com/julienschmidt/httprouter"
)

var methods = map[string]string{
	"Get":     http.MethodGet,
	"Put":     http.MethodPut,
	"Post":    http.MethodPost,
	"Head":    http.MethodHead,
	"Patch":   http.MethodPatch,
	"Delete":  http.MethodDelete,
	"Connect": http.MethodConnect,
	"Options": http.MethodOptions,
	"Trace":   http.MethodTrace,
}

// Middleware server middleware
type Middleware func(resp http.ResponseWriter, req *http.Request, next http.Handler)

// Server rpc server
type Server interface {
	Handle(path string, service interface{}, middleware ...Middleware) Server
	ServeHTTP(w http.ResponseWriter, r *http.Request)
}

type middlewareHandler struct {
	middleware Middleware
	next       http.Handler
}

func (handler *middlewareHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	handler.middleware(w, r, handler.next)
}

type serverImpl struct {
	slf4go.Logger
	router *httprouter.Router
}

// New create new Server
func New() Server {
	return &serverImpl{
		Logger: slf4go.Get("server"),
		router: httprouter.New(),
	}
}

func (server *serverImpl) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	server.router.ServeHTTP(w, r)
}

func (server *serverImpl) checkInputType(paramT reflect.Type) bool {
	if paramT.Kind() != reflect.Ptr {
		return false
	}

	if paramT.Elem().Kind() != reflect.Struct {
		return false
	}

	return true
}

func (server *serverImpl) Handle(path string, service interface{}, middleware ...Middleware) Server {

	serviceT := reflect.TypeOf(service)

	for i := 0; i < serviceT.NumMethod(); i++ {
		method := serviceT.Method(i)
		println(method.Name)

		httpMethod, ok := server.validMethod(method.Name)

		if !ok {
			server.DebugF("[%s] skip invalid method %s", serviceT, method.Name)
			continue
		}

		if method.Type.NumIn() != 3 {
			server.DebugF("[%s] skip invalid method %s,input parameters != 2 (%d)", serviceT, method.Name, method.Type.NumIn())
			continue
		}

		if !server.checkInputType(method.Type.In(1)) {
			server.DebugF("[%s] skip invalid method %s param 1 %s , parameter must be struct ptr", serviceT, method.Name, method.Type.In(1))
			continue
		}

		if !server.checkInputType(method.Type.In(2)) {
			server.DebugF("[%s] skip invalid method %s param 2 %s , parameter must be struct ptr", serviceT, method.Name, method.Type.In(2))
			continue
		}

		if method.Type.NumOut() != 1 {
			server.DebugF("[%s] skip invalid method %s,output parameters != 1 ", serviceT, method.Name)
			continue
		}

		var err error

		if !method.Type.Out(0).Implements(reflect.TypeOf(&err).Elem()) {
			server.DebugF("[%s] skip invalid method %s,out parameter not implement error interface", serviceT, method.Name)
			continue
		}

		handler := server.packageHandlers(server.createHandle(service, &method), middleware...)

		server.router.Handler(httpMethod, path, handler)

		server.InfoF("[%s] find valid http %s method %s ", serviceT, httpMethod, method.Name)
	}

	return server
}

func (server *serverImpl) packageHandlers(handler http.Handler, middlewares ...Middleware) http.Handler {

	next := handler

	for i := len(middlewares); i > 0; i-- {
		next = &middlewareHandler{
			middleware: middlewares[i-1],
			next:       next,
		}
	}

	return next
}

func (server *serverImpl) validMethod(name string) (string, bool) {
	for prefix, method := range methods {
		if strings.HasPrefix(name, prefix) {
			return method, true
		}
	}

	return "", false
}

func (server *serverImpl) createHandle(service interface{}, method *reflect.Method) http.Handler {

	serviceValue := reflect.ValueOf(service)

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		input, err := server.readParameter(r, method.Type.In(1))

		if err != nil {
			server.writeResponse(w, nil, err)
			return
		}

		params := []reflect.Value{
			serviceValue,
			input,
			reflect.New(method.Type.In(2)),
		}

		results := method.Func.Call(params)

		err, ok := results[0].Interface().(error)

		if !ok {
			panic(fmt.Sprintf("filter service %s RESTful method %s error,result must be error", reflect.TypeOf(service), method.Name))
		}

		server.writeResponse(w, results[0].Interface(), err)
	})
}

func (server *serverImpl) writeResponse(w http.ResponseWriter, result interface{}, err error) {

}

func (server *serverImpl) readParameter(r *http.Request, paramT reflect.Type) (reflect.Value, error) {

	return reflect.New(paramT), nil
}
