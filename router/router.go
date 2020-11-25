package router

import (
    "GoChessgameServer/controllers"

    "github.com/gorilla/mux"
)

// API Router of the program
var Router *mux.Router

func init() {
    // Create our API router
    Router = mux.NewRouter()

    // Handlers on a variety of controllers
    Router.HandleFunc("/api/user/create", controllers.CreateLogin).Methods("POST")
    Router.HandleFunc("/api/user/login", controllers.LoginUsers).Methods("POST")
    Router.HandleFunc("/api/hello", controllers.HelloController).Methods("GET")

    // Middleware for token checking
    Router.Use(controllers.TokenChecker)
}
