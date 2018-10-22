package validator

import (
	"sync"

	"github.com/dynamicgo/slf4go"
)

var logger slf4go.Logger
var once sync.Once

func initLogger() {
	once.Do(func() {
		logger = slf4go.Get("validator")
	})
}
