package controllers

import (
    "net/http"
    "encoding/json"
    "io/ioutil"
    "reflect"

    u "GoChessgameServer/util"
    "GoChessgameServer/database"
)

// This function returns an array of user statistic
// It can be used, for example, for displaying user records table
func UserStatistic(w http.ResponseWriter, r *http.Request) {

    writeError := u.WriteErrorCreator(w)

    contextUser := r.Context().Value("login").(string)
    var results *database.RowsResult
    var err error

    // Case when is no body
    if r.Body == http.NoBody {
        // If it is no body, we will return statistic of all players
        results, err = database.QueryBlocking(`
        SELECT
            U.Login AS Login,
            COUNT(G.ID) AS GamesPlayed,
            SUM(CASE WHEN G.WinnerLogin = U.Login THEN 1 ELSE 0 END) AS WinnedGames
        FROM
            dbo.Users U LEFT JOIN dbo.GamesHistory G ON G.PlayerOneLogin = U.Login OR G.PlayerTwoLogin = U.Login
        GROUP BY
            U.Login
        `)

    } else {
        jsonString, err := ioutil.ReadAll(r.Body)
        if err != nil {
            writeError("Connection error")
            w.WriteHeader(http.StatusInternalServerError)
            contrLogger.Printf("UserStatistic: Error when reading user request: %s\n", err.Error())
            return
        }

        // Try to unmarshall to self request
        passed := true

        type reqType struct {
            Self bool `json:"self"`
        }
        req := reqType{}

        err = json.Unmarshal(jsonString, &req)
        if err != nil || reflect.DeepEqual(req, reqType{}) {
            // Its not self request, try to parse next request type
            passed = false
        }

        if passed {
            // Make database query for self request

            results, err = database.QueryBlocking(`
            SELECT
                U.Login AS Login,
                COUNT(G.ID) AS GamesPlayed,
                SUM(CASE WHEN G.WinnerLogin = U.Login THEN 1 ELSE 0 END) AS WinnedGames
            FROM
                dbo.Users U LEFT JOIN dbo.GamesHistory G ON G.PlayerOneLogin = U.Login OR G.PlayerTwoLogin = U.Login
            WHERE
                U.Login = $1
            GROUP BY
                U.Login
            `, contextUser)
        } else {
            // Try to unmarshal user request
            type reqType struct {
                Login string `json:"login"`
            }
            req := reqType{}

            err = json.Unmarshal(jsonString, &req)
            if err != nil || reflect.DeepEqual(req, reqType{}) {
                // Its also not user request, invalid JSON request
                writeError("Invalid/Mailformed JSON request")
                return
            }

            // User validation
            if !u.ValidateCredentials(req.Login, "", "") {
                writeError("Passed login is not corresponds to login requirements")
                return
            }

            // Make database query for user request
            results, err = database.QueryBlocking(`
            SELECT
                U.Login AS Login,
                COUNT(G.ID) AS GamesPlayed,
                SUM(CASE WHEN G.WinnerLogin = U.Login THEN 1 ELSE 0 END) AS WinnedGames
            FROM
                dbo.Users U LEFT JOIN dbo.GamesHistory G ON G.PlayerOneLogin = U.Login OR G.PlayerTwoLogin = U.Login
            WHERE
                U.Login = $1
            GROUP BY
                U.Login
            `, req.Login)
        }
    }

    if err != nil {
        writeError("Connection error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("UserStatistic: Error when executing query: %s\n", err.Error())
        return
    }
    if len(*results) == 0 {
        writeError("Oops, it seems that you account has been deleted. Please, restart you application")
        return
    }

    // Parse results
    type respRow struct {
        User string `json:"user"`
        GamesPlayed int64 `json:"played"`
        WinnedGames int64 `json:"winned"`
    }
    resp := []respRow{}

    for _, row := range *results {
        resp = append(resp, respRow{
            User: row["Login"].(string),
            GamesPlayed: row["GamesPlayed"].(int64),
            WinnedGames: row["WinnedGames"].(int64),
        })
    }

    // Send response
    if err = json.NewEncoder(w).Encode(resp); err != nil {
        writeError("Server error")
        w.WriteHeader(http.StatusInternalServerError)
        contrLogger.Printf("IsAdmin: Error when sending response: %s\n", err.Error())
        return
    }
    w.Header().Add("Content-Type", "application/json")

    // Log new user
    contrLogger.Printf("UserStatistic: User %s requested user statistic\n", contextUser)
}
