package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

type RequestHandler func (http.ResponseWriter, *http.Request)
type Middleware func (*RequestHandler) RequestHandler

func WithLogger(next *RequestHandler) RequestHandler {
	return func (res http.ResponseWriter, req *http.Request)  {
		fmt.Printf("--> %s %s\n", req.Method, req.URL)
		start := time.Now()
        if next != nil {
            (*next)(res, req)
        }
		fmt.Printf("<-- %s %s %dms\n", req.Method, req.URL, time.Since(start).Milliseconds())
	}
}
