package main

import (
	"ingrid/pkg/webapp"
	"net/http"
)


func main() {
	http.HandleFunc("/routes", webapp.HandleRequest)
	http.ListenAndServe(":8080", nil)
}
