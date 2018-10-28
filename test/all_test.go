package test

import (
	"encoding/json"
	"net/http"
	"testing"

	"github.com/dynamicgo/restrpc/server"
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

func init() {
	server := server.New()
	server.Handle("/", &A{})
	http.ListenAndServe(":8080", server)
}

func TestHandle(t *testing.T) {

}

func printResult(v interface{}) string {
	val, _ := json.MarshalIndent(v, "", "\t")

	return string(val)
}
