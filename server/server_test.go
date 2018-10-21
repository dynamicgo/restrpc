package server

import (
	"encoding/json"
	"testing"
)

type A struct {
}

type Param struct {
}

type Result struct {
}

func (a *A) GetMessage(p *Param, r *Result) error {
	return nil
}

func (a *A) PostMessage(p *Param, r *Result) error {
	return nil
}

func (a *A) PostB(p *Param, r *Result) error {
	return nil
}

func TestHandle(t *testing.T) {
	server := New()
	server.Handle("/", &A{})
}

func printResult(v interface{}) string {
	val, _ := json.MarshalIndent(v, "", "\t")

	return string(val)
}
