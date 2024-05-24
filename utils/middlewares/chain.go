package middlewares

import (
	"slices"

	"github.com/rohhamh/go-shopping-cart-crud/middlewares"
)

func Chain (middlewaresChain *[]middlewares.Middleware, handler *middlewares.RequestHandler) middlewares.RequestHandler {
    if middlewaresChain == nil { return *handler }

    slices.Reverse(*middlewaresChain)
    for _, middleware := range *middlewaresChain {
        middlewareHandler := middleware(handler)
        handler = &middlewareHandler
    }
    slices.Reverse(*middlewaresChain)
    return *handler
}
