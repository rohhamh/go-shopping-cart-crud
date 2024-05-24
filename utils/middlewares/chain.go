package middlewares

import (
	"reflect"
	"runtime"
	"slices"

	"github.com/rohhamh/go-shopping-cart-crud/middlewares"
)

func GetFunctionName(i interface{}) string {
    return runtime.FuncForPC(reflect.ValueOf(i).Pointer()).Name()
}

func Chain (middlewaresChain *[]middlewares.Middleware, handler *middlewares.RequestHandler) middlewares.RequestHandler {
    if middlewaresChain == nil { return *handler }

    slices.Reverse(*middlewaresChain)
    middlewaresHandlers := []middlewares.RequestHandler {*handler}
    finalHandler := handler
    for _, middleware := range *middlewaresChain {
        middlewareHandler := middleware(finalHandler)
        finalHandler = &middlewareHandler
        middlewaresHandlers = append(middlewaresHandlers, *finalHandler)
    }
    return *finalHandler
}
