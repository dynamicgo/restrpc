package server

import (
	"encoding/json"
	"testing"
)

type A struct {
}

type param struct {
}

type result struct {
}

func (a *A) GetMessage(p *param, r *result) error {
	return nil
}

func (a *A) PostMessage(p *param, r *result) error {
	return nil
}

func (a *A) PostB(p *param, r *result) error {
	return nil
}

func TestHandle(t *testing.T) {
	server := New()
	server.Handle("/", &A{})

	var c chan int

	<-c
}

func printResult(v interface{}) string {
	val, _ := json.MarshalIndent(v, "", "\t")

	return string(val)
}
