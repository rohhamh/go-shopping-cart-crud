package middlewares

import (
	"fmt"
	"net/http"
	"reflect"
	"runtime"
)

func GetFunctionName(i interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func Authorize(next *RequestHandler) RequestHandler {
	return func (res http.ResponseWriter, req *http.Request)  {
        fmt.Printf("headers %v\n", req.Header)
        if next != nil {
            (*next)(res, req)
        }
	}
}
