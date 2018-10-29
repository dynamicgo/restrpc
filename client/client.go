package client

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"path/filepath"
	"strings"

	"github.com/dynamicgo/xerrors/apierr"

	"github.com/dynamicgo/restrpc"
	"github.com/dynamicgo/xerrors"
	"github.com/go-resty/resty"
)

// Auth .
type Auth interface {
	Handle(request *http.Request)
}

// Option .
type Option func(request *http.Request)

// WithAuth add auth option
func WithAuth(auth Auth) Option {
	return func(request *http.Request) {
		auth.Handle(request)
	}
}

// WithJWToken .
func WithJWToken(token string) Option {
	return func(request *http.Request) {
		request.Header.Add("Authorization", fmt.Sprintf("Bearer %s", token))
	}
}

// ErrType .
var (
	ErrMethod = errors.New("unsupport method")
)

// Client .
type Client interface {
	Call(path string, method string, args interface{}, reply interface{}, options ...Option) error
	Service(path string) Service
}

// Service .
type Service interface {
	Call(method string, name string, args interface{}, reply interface{}, options ...Option) error
}

type clientImpl struct {
	rootURL string
}

// New .
func New(url string) Client {
	return &clientImpl{
		rootURL: url,
	}
}

func (client *clientImpl) Call(path string, method string, args interface{}, reply interface{}, options ...Option) error {

	pathnodes := strings.Split(path, "/")

	name := pathnodes[len(pathnodes)-1]

	path = strings.Join(pathnodes[:len(pathnodes)-1], "/")

	return client.Service(path).Call(method, name, args, reply, options...)

}

type serviceImpl struct {
	rootURL string
	path    string
}

func (client *clientImpl) Service(path string) Service {
	return &serviceImpl{
		rootURL: client.rootURL,
		path:    path,
	}
}

func (service *serviceImpl) Call(method string, name string, args interface{}, reply interface{}, options ...Option) error {

	url := fmt.Sprintf("%s/%s/%s", service.rootURL, service.path, name)

	checkedURL, err := service.checkURL(url)

	if err != nil {
		return xerrors.Wrapf(err, "check url %s failed", url)
	}

	switch method {
	case http.MethodGet:
		return service.Get(checkedURL, args, reply, options...)
	case http.MethodPost:
		return service.Post(checkedURL, args, reply, options...)
	case http.MethodDelete:
		return service.Delete(checkedURL, args, reply, options...)
	case http.MethodPut:
		return service.Put(checkedURL, args, reply, options...)
	default:
		return xerrors.Wrapf(ErrMethod, "invalid method %s", method)
	}

}

func (service *serviceImpl) Get(checkedURL string, args interface{}, reply interface{}, options ...Option) error {
	r := resty.R().SetQueryParams(service.args2Map(args)).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	for _, option := range options {
		option(r.RawRequest)
	}

	resp, err := r.Get(checkedURL)

	if err != nil {
		return xerrors.Wrapf(err, "network error")
	}

	return service.checkResult(resp, reply)
}

func (service *serviceImpl) Post(checkedURL string, args interface{}, reply interface{}, options ...Option) error {
	r := resty.R().SetBody(args).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	for _, option := range options {
		option(r.RawRequest)
	}

	resp, err := r.Post(checkedURL)

	if err != nil {
		return xerrors.Wrapf(err, "network error")
	}

	return service.checkResult(resp, reply)
}

func (service *serviceImpl) Delete(checkedURL string, args interface{}, reply interface{}, options ...Option) error {
	r := resty.R().SetQueryParams(service.args2Map(args)).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	for _, option := range options {
		option(r.RawRequest)
	}

	resp, err := r.Delete(checkedURL)

	if err != nil {
		return xerrors.Wrapf(err, "network error")
	}

	return service.checkResult(resp, reply)
}

func (service *serviceImpl) Put(checkedURL string, args interface{}, reply interface{}, options ...Option) error {
	r := resty.R().SetBody(args).
		SetHeader("Content-Type", "application/json").
		SetHeader("Accept", "application/json")

	for _, option := range options {
		option(r.RawRequest)
	}

	resp, err := r.Put(checkedURL)

	if err != nil {
		return xerrors.Wrapf(err, "network error")
	}

	return service.checkResult(resp, reply)
}

type result struct {
	Code   int         `json:"code"`
	ErrMsg string      `json:"errmsg"`
	Result interface{} `json:"result"`
}

func (service *serviceImpl) checkResult(resp *resty.Response, reply interface{}) error {

	var r result

	err := json.Unmarshal(resp.Body(), &r)

	if err != nil {
		return xerrors.Wrapf(restrpc.ErrInternal, "unmarshal %s err %s", resp.Body(), err)
	}

	if resp.StatusCode() != http.StatusOK {
		return xerrors.Wrapf(apierr.New(r.Code, r.ErrMsg), "apierr: %s", resp.Body())
	}

	if r.Result == nil {
		return xerrors.Wrapf(restrpc.ErrInternal, "result not found: %s", resp.Body())
	}

	buff, err := json.Marshal(r.Result)

	if err != nil {
		return xerrors.Wrapf(restrpc.ErrInternal, "marshal result %v err %s", r.Result, err)
	}

	if err := json.Unmarshal(buff, reply); err != nil {
		return xerrors.Wrapf(restrpc.ErrInternal, "unmarshal %s err %s", resp.Body(), err)
	}

	return nil
}

func (service *serviceImpl) checkURL(s string) (string, error) {
	u, err := url.Parse(s)

	if err != nil {
		return "", err
	}

	u.Path = filepath.Clean(u.Path)

	return u.String(), nil
}

func (service *serviceImpl) args2Map(args interface{}) map[string]string {
	return nil
}
