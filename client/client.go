package client

import (
	"errors"
)

// ErrType .
var (
	ErrType = errors.New("invalid client interface type")
)

// Client .
type Client interface {
	Call(path string, args interface{}, reply interface{}) error
	Service(path string) Service
}

// Service .
type Service interface {
	Call(method string, args interface{}, reply interface{}) error
}
