package clientapi

import (
    "net/http"
    "encoding/json"

    u "GoChessgameServer/util"
    "GoChessgameServer/database"
)

// This function returns an array of user statistic
// It can be used, for example, for displaying user records table
func UserStatistic(w http.ResponseWriter, r *http.Request) {

    contextUser := r.Context().Value("login").(string)
    writeError := u.WriteErrorCreator(w)

    // Type for statistic values
    type userStatistic struct {
        User string `json:"user"`
        GamesPlayed int64 `json:"played"`
        WinnedGames int64 `json:"winned"`
    }

    userStatistics := []userStatistic{}

    var err error
    baseQuery := database.DB.
        Table("users").
        Joins("JOIN games_history ON games_history.player_one_login = users.login OR games_history.player_two_login = users.login").
        Select(`
            users.login as user,
            count(games_history.id) as games_played,
            sum(case where games_history.winner_login = users.login then 1 else 0 end) as winned_games
        `).
        Group("users.login")

    if r.Body == http.NoBody {
        // Query all users from database
        err = baseQuery.
            Find(&userStatistics).Error
    } else {
        // Parse json request at first
        response := struct{
            Login string `json:"login"`
        }{}

        if json.NewDecoder(r.Body).Decode(&response); err != nil {
            writeError("Invalid parameters")
            clientApiLogger.Printf("UserStatistic: Error when reading user request: %s\n", err.Error())
            return
        }

        // Validate login
        success := u.ValidateValues(
            &u.VValue{Type: "Login", Value: response.Login},
        )

        if !success {
            writeError("Passed login doesn't satisfy value requirements")
            return
        }

        // Query information with user passed filter
        err = baseQuery.
            Where("users.login = ?", response.Login).
            Find(&userStatistics).Error
    }

    if err != nil {
        writeError("Connection error")
        w.WriteHeader(http.StatusInternalServerError)
        clientApiLogger.Printf("UserStatistic: Error when executing query: %s\n", err.Error())
        return
    }

    // Send response
    if err = json.NewEncoder(w).Encode(userStatistics); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        clientApiLogger.Printf("IsAdmin: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")

    // Log
    clientApiLogger.Printf("UserStatistic: User %s requested user statistic\n", contextUser)
}
