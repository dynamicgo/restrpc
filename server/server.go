package server

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"

	"github.com/dynamicgo/xerrors/apierr"

	"github.com/dynamicgo/xerrors"

	"github.com/dynamicgo/restrpc/validator"

	"github.com/dynamicgo/slf4go"

	"github.com/julienschmidt/httprouter"
)

// Errors
var (
	ErrUnsupportContentType = errors.New("unsupport content-type")
	ErrInternal             = apierr.New(-1, "INNER_ERROR")
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

// R Response type
type R map[string]interface{}

// Middleware server middleware
type Middleware func(resp http.ResponseWriter, req *http.Request, next http.Handler)

// Server rpc server
type Server interface {
	Handle(path string, service interface{}, middleware ...Middleware) Server
	ServeHTTP(w http.ResponseWriter, r *http.Request)
	Success(w http.ResponseWriter, result interface{}) error
	Fail(w http.ResponseWriter, code int, cause error) error
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

		name := strings.TrimPrefix(strings.ToLower(method.Name), strings.ToLower(httpMethod))

		server.router.Handler(httpMethod, fmt.Sprintf("%s/%s", path, name), handler)

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
			server.writeResponse(w, nil, http.StatusInternalServerError, err)
			return
		}

		output := reflect.New(method.Type.In(2))

		params := []reflect.Value{
			serviceValue,
			input,
			output,
		}

		results := method.Func.Call(params)

		err, ok := results[0].Interface().(error)

		if !ok {
			panic(fmt.Sprintf("filter service %s RESTful method %s error,result must be error", reflect.TypeOf(service), method.Name))
		}

		if err != nil {
			server.writeResponse(w, nil, http.StatusInternalServerError, err)
		} else {
			server.writeResponse(w, R{
				"result": output.Interface(),
			}, http.StatusOK, nil)
		}

	})
}

func (server *serverImpl) writeResponse(w http.ResponseWriter, r R, code int, err error) error {

	if err != nil {
		apiErr := apierr.As(err, ErrInternal)
		r["code"] = apiErr.Code()
		r["errmsg"] = apiErr.Error()
	}

	buff, err := json.Marshal(r)

	if err != nil {
		return xerrors.Wrapf(err, "marshal response %v err %s", r, err)
	}

	w.WriteHeader(code)
	_, err = w.Write(buff)

	if err != nil {
		return xerrors.Wrapf(err, "write response err")
	}

	return nil
}

func (server *serverImpl) Success(w http.ResponseWriter, result interface{}) error {
	return server.writeResponse(w, R{
		"result": result,
	}, http.StatusOK, nil)
}
func (server *serverImpl) Fail(w http.ResponseWriter, code int, cause error) error {
	return server.writeResponse(w, nil, code, cause)
}

func (server *serverImpl) readParameter(r *http.Request, paramT reflect.Type) (reflect.Value, error) {

	var reader validator.Reader

	if r.Method == http.MethodGet || r.Method == http.MethodDelete {
		reader = validator.NewQueryReader(r.URL.Query())
	} else {

		contentType := r.Header.Get("Content-type")

		if contentType != "application/json" {
			return reflect.Value{}, xerrors.Wrapf(ErrUnsupportContentType, "restrpc only support application/json content-type")
		}

		buff, err := ioutil.ReadAll(r.Body)

		if err != nil {
			return reflect.Value{}, xerrors.Wrapf(err, "unable read request body from %s", r.RequestURI)
		}

		reader, err = validator.NewJSONReader(buff)

		if err != nil {
			return reflect.Value{}, xerrors.Wrapf(err, "parse request %s json body error", r.RequestURI)
		}
	}

	values, err := validator.Validate(reader, paramT)

	if err != nil {
		return reflect.Value{}, err
	}

	return values[0], nil
}
