package longpollmanagers

import (
    "net/http"
    "strconv"

    u "GoChessgameServer/util"
    "GoChessgameServer/store"

    lp "github.com/jcuga/golongpoll"
)

// This function returns function handler,
// that controls end game event
func EndGame(manager *lp.LongpollManager) func(http.ResponseWriter, *http.Request) {

    return func(w http.ResponseWriter, r *http.Request) {

        category := r.URL.Query().Get("category")
        contextUser := r.Context().Value("login").(string)
        writeError := u.WriteErrorCreator(w)

        // Parse game id
        gameId, err := strconv.Atoi(category)
        if err != nil {
            writeError("Passed gameId is invalid, parse valid int value")
            return
        }

        // Get game store
        gameStore, err := store.GetGameStore(gameId)
        if err != nil {
            writeError("Game with the specified is not found, aborting")
            return
        }

        // Check game started flag
        if !gameStore.GameStarted {
            writeError("Cant end game, game is not started")
            return
        }

        // Check if its a user, that really need to know about game ending
        if contextUser != gameStore.PlayerOneLogin && contextUser != gameStore.PlayerTwoLogin {
            writeError("You are not permitted to end this game, aborting")
            lmLogger.Printf("EndGame: User %s tryied to end other game\n", contextUser)
            return
        }

        // Perform game end subscription
        manager.SubscriptionHandler(w, r)
    }
}
