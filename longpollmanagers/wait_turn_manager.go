package longpollmanagers

import (
    "net/http"
    "strconv"

    u "GoChessgameServer/util"
    "GoChessgameServer/store"

    lp "github.com/jcuga/golongpoll"
)

// This function return functino handler,
// that allows players to listen opponent's turn
func WaitTurn(manager *lp.LongpollManager) func(http.ResponseWriter, *http.Request) {

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
            writeError("Can't send turn, game is not started")
            return
        }

        // Check for player turn
        if gameStore.IsPlayerOneTurn && contextUser != gameStore.PlayerTwoLogin || !gameStore.IsPlayerOneTurn && contextUser != gameStore.PlayerOneLogin {
            writeError("You are not permitted to make turn. There are opponent turn right now")
            lmLogger.Printf("WaitTurn: Player %s tried to use his turn when its opponent turn in game id %d\n", contextUser, gameId)
            return
        }

        // Perform game end subscription
        manager.SubscriptionHandler(w, r)
    }
}
