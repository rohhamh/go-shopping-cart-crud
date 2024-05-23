package handlers

import (
	"fmt"
	"net/http"
	"time"
)

type requestHandler func (http.ResponseWriter, *http.Request)

func WithLogger(httpHandler requestHandler) requestHandler {
	return func (res http.ResponseWriter, req *http.Request)  {
		fmt.Printf("--> %s %s\n", req.Method, req.URL)
		start := time.Now()
		httpHandler(res, req)
		fmt.Printf("<-- %s %s %dms\n", req.Method, req.URL, time.Since(start).Milliseconds())
	}
}
