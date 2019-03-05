package controllers

import (
	"fmt"
	"log"
	"net/http"
)

// SayHello say hello and show source address
func SayHello(w http.ResponseWriter, req *http.Request) {
	log.Println("Source:", req.RemoteAddr, "Request URI:", req.RequestURI)
	fmt.Fprintf(w, "Hello!\nYour network address is: %v", req.RemoteAddr)
}
