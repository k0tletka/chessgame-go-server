package router

import (
    "GoChessgameServer/controllers"
    lm "GoChessgameServer/longpollmanagers"
    "GoChessgameServer/store"

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
    Router.HandleFunc("/api/user/info", controllers.UserInfo).Methods("GET")
    Router.HandleFunc("/api/lobby/list", controllers.LobbyList).Methods("GET")
    Router.HandleFunc("/api/lobby/create", controllers.LobbyCreate).Methods("POST")
    Router.HandleFunc("/api/game/connect", controllers.GameConnect).Methods("POST")
    Router.HandleFunc("/api/game/ack", controllers.GameAck).Methods("POST")
    Router.HandleFunc("/api/playerstat", controllers.UserStatistic).Methods("GET", "POST")
    Router.HandleFunc("/api/motd", controllers.GetMotd).Methods("GET")

    // Longpoll managers
    Router.HandleFunc("/api/game/wait", lm.WaitUser(store.WaitGameLM)).Methods("GET")

    // Middleware for token checking
    Router.Use(controllers.TokenChecker)
}
