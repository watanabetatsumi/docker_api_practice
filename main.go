package main

import (
	"fmt"
	"net/http"
	
)

func main() {
    http.HandleFunc("/", restHandler)
    http.ListenAndServe(":8080", nil)
}

func restHandler(w http.ResponseWriter, r *http.Request) {
    if r.Method == "GET" {
        fmt.Fprintln(w, "GET called!!")
    } else if r.Method == "POST" {
        fmt.Fprintln(w, "POST called!!")
    } else if r.Method == "PUT" {
        fmt.Fprintln(w, "PUT called!!")
    } else if r.Method == "DELETE" {
        fmt.Fprintln(w, "DELETE called!!")
    }
}
