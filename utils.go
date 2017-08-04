package axe

import (
	"reflect"
	"runtime"
)

// FuncName get func name
func FuncName(f interface{}) string {
	return runtime.FuncForPC(reflect.ValueOf(f).Pointer()).Name()
}
