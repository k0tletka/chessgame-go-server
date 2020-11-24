package controllers

import (
    "net/http"
    "io"
    "fmt"
)

// Sample controller for API testing
func HelloController(w http.ResponseWriter, r *http.Request) {
    contrLogger.Println(fmt.Sprintf("Request on %s", r.URL.Path))
    io.WriteString(w, "Hello, world!")
}
