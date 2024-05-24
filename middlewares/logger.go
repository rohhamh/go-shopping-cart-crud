package middlewares

import (
	"fmt"
	"net/http"
	"time"
)

type RequestHandler func (http.ResponseWriter, *http.Request)
type Middleware func (*RequestHandler) RequestHandler

type ResponseWriter struct {
    http.ResponseWriter
    statusCode              int
}

func (rwsc *ResponseWriter) WriteHeader(code int) {
    rwsc.statusCode = code
    rwsc.ResponseWriter.WriteHeader(code)
}

func Logger(next *RequestHandler) RequestHandler {
	return func (res http.ResponseWriter, req *http.Request)  {
		fmt.Printf("--> %s %s\n", req.Method, req.URL)
		start := time.Now()
        responseWriter := &ResponseWriter { res, 200 }
        if next != nil {
            (*next)(responseWriter, req)
        }
		fmt.Printf("<-- %s %s %v %dms\n",
            req.Method, req.URL, responseWriter.statusCode,
            time.Since(start).Milliseconds())
	}
}
