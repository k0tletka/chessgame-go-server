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
    Router.HandleFunc("/api/user/isadmin", controllers.IsAdmin).Methods("GET")
    Router.HandleFunc("/api/user/changepass", controllers.ChangePassword).Methods("POST")
    Router.HandleFunc("/api/playerstat", controllers.UserStatistic).Methods("GET", "POST")

    // Middleware for token checking
    Router.Use(controllers.TokenChecker)
}
