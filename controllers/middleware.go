package controllers

import (
	"log"
	"net/http"
)

// MiddleWareLogOnCalled a middleware to log request info
func MiddleWareLogOnCalled(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
		log.Println("Request:", req.Method, req.RequestURI)
		log.Println("From:", req.RemoteAddr)
		next.ServeHTTP(w, req)
	})
}
