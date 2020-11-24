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
    Router.HandleFunc("/hello", controllers.HelloController).Methods("GET")
}
