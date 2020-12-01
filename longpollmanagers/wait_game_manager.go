package longpollmanagers

import (
    "net/http"
    "strconv"

    u "GoChessgameServer/util"
    "GoChessgameServer/store"

    lp "github.com/jcuga/golongpoll"
)

// This function returns a new http handler,
// that controls user waiting
func WaitUser(manager *lp.LongpollManager) func(http.ResponseWriter, *http.Request) {

    return func(w http.ResponseWriter, r *http.Request) {

        category := r.URL.Query().Get("category")
        contextUser := r.Context().Value("login")
        writeError := u.WriteErrorCreator(w)

        // Parse game id
        gameId, err := strconv.Atoi(category)
        if err != nil {
            writeError("Passed gameId is invalid, parse valid int value")
            return
        }

        // Get game store with passed gameId
        store, err := store.GetGameStore(gameId)
        if err != nil {
            writeError("Game with the specified is not found, aborting")
            return
        }

        // Check if its a user, that really need to wait an opponent
        if contextUser != store.PlayerOneLogin {
            writeError("You are not permitted to wait for this game, aborting")
            lmLogger.Printf("WaitUser: User %s tryied to wait for other game\n", contextUser)
            return
        }

        // All ok, start longpoll handle
        manager.SubscriptionHandler(w, r)
    }
}
