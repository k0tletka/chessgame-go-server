package game

import (
    "encoding/json"
    "time"

    "GoChessgameServer/database"
    u "GoChessgameServer/util"

    "github.com/gorilla/websocket"
)

// This type represents json model for each turn, that clients making
type Turn struct {
    GameID      int     `json:"game_id"`
    Login       string  `json:"login"`
    FigposX     int     `json:"figposx"`
    FigposY     int     `json:"figposy"`
    AltX        int     `json:"xalt"`
    AltY        int     `json:"yalt"`
    Surrender   bool    `json:"surrender"`
}

// This type represents end game json model
type endGameResult struct {
    IsDraw      bool    `json:"is_draw"`
    WinnerLogin string  `json:"winner_login"`
}

// This function is a main entrypoint handler,
// that controls game process.
// This function starts as coroutine
func controlGame(gameSession *GameSession) {

    // Fill table with figures
    table := createNewTable()

    // Turner to get rid about who must make turn next
    var turner *GameTurner

    gameStartedTimestamp := time.Now()
    whiteTurn := true
    performNextTurn := true

    // Get copy of players
    players := []string{}
    for _, v := range gameSession.GetAllPlayers() {
        players = append(players, v.Login)
    }

    // Handle turns
    turnTimer := time.NewTimer(time.Minute * 5)

    for {
        // Get player that must make next turn
        if performNextTurn {
            if turner == nil {
                turner = gameSession.PlayerTurner
            } else {
                turner = turner.Next
            }

            performNextTurn = false
        }

        conn := turner.GameConnection
        wsConn := conn.Connection.GetConnection()

        select {
        case <-turnTimer.C:
            broadcastJSONMessages(&endGameResult{IsDraw: true}, conn, gameSession)

            // Turn time was exceed
            saveResultsToDatabase(gameSession, players, true, "", gameStartedTimestamp)
            closeConnections(gameSession)
            return
        case turn := <-conn.ReadChannel:
            if turn.Surrender {
                // Close connection for client that surrended
                conn.Connection.CloseConnection("You have surrended. Game ended")
                gameSession.deleteConnection(conn)

                // Check that at least two connections left. If amount of connections
                // is lower that two, chose winner and close session
                if aPl := gameSession.GetAllPlayers(); len(aPl) == 1 {
                    broadcastJSONMessages(&endGameResult{WinnerLogin: aPl[0].Login}, conn, gameSession)
                    saveResultsToDatabase(gameSession, players, false, aPl[0].Login, gameStartedTimestamp)
                    closeConnections(gameSession)
                    return
                } else if len(aPl) < 1 {
                    return
                }

                performNextTurn = true
            }

            // Check figure existent
            if table[turn.FigposX][turn.FigposY] == nil {
                wsConn.WriteMessage(websocket.TextMessage, u.ErrorJson("There is no figure"))
                continue
            }

            // Check for player turn and perform operation
            if (whiteTurn && table[turn.FigposX][turn.FigposY].IsFigureBlack()) ||
                (!whiteTurn && !table[turn.FigposX][turn.FigposY].IsFigureBlack()) {

                wsConn.WriteMessage(websocket.TextMessage, u.ErrorJson("You can't move figure that you don't own"))
                continue
            }

            if !table[turn.FigposX][turn.FigposY].CanFigurePass(turn.AltX, turn.AltY) {
                wsConn.WriteMessage(websocket.TextMessage, u.ErrorJson("You can't move this figure at the specified location"))
                continue
            }

            // Pass table[turn.FigposX][turn.FigposY]
            _ = table[turn.FigposX][turn.FigposY].Pass(turn.AltX, turn.AltY)

            // Check win state
            if checkBlackWin(table) {
                var winnerGameObject *GameClientConnection

                if !whiteTurn {
                    winnerGameObject = conn
                } else {
                    turner = turner.Next
                    if turner == nil {
                        winnerGameObject = gameSession.PlayerTurner.GameConnection
                    } else { winnerGameObject = turner.GameConnection }
                }

                broadcastJSONMessages(&endGameResult{WinnerLogin: winnerGameObject.Login}, conn, gameSession)

                // Write results to database
                saveResultsToDatabase(gameSession, players, false, winnerGameObject.Login, gameStartedTimestamp)
                closeConnections(gameSession)
                return
            }

            if checkWhiteWin(table) {
                var winnerGameObject *GameClientConnection

                if whiteTurn {
                    winnerGameObject = conn
                } else {
                    turner = turner.Next
                    if turner == nil {
                        winnerGameObject = gameSession.PlayerTurner.GameConnection
                    } else { winnerGameObject = turner.GameConnection }
                }

                broadcastJSONMessages(&endGameResult{WinnerLogin: winnerGameObject.Login}, conn, gameSession)

                // Write results to database
                saveResultsToDatabase(gameSession, players, false, winnerGameObject.Login, gameStartedTimestamp)
                closeConnections(gameSession)
                return
            }

            performNextTurn = true
            whiteTurn = !whiteTurn

            broadcastJSONMessages(&turn, conn, gameSession)
        case <-conn.Connection.ConnectionClosed:
            gameSession.deleteConnection(conn)

            // Check that at least two connections left. If amount of connections
            // is lower that two, chose winner and close session
            if aPl := gameSession.GetAllPlayers(); len(aPl) == 1 {
                saveResultsToDatabase(gameSession, players, false, aPl[0].Login, gameStartedTimestamp)
                closeConnections(gameSession)
                return
            } else if len(aPl) < 1 {
                return
            }
        }

    }
}

// This function creates new table with filled figures in stardart places
func createNewTable() *ChessTable {
    table := &ChessTable{}

    // Pawn
    for i := 0; i <= 7; i++ {
        table[1][i] = &Pawn{cTable: table, x: 1, y: i, isBlack: false, isMoved: false}
        table[6][i] = &Pawn{cTable: table, x: 6, y: i, isBlack: true, isMoved: false}
    }

    // Rook
    table[0][0] = &Rook{cTable: table, x: 0, y: 0, isBlack: false}
    table[0][7] = &Rook{cTable: table, x: 0, y: 7, isBlack: false}
    table[7][0] = &Rook{cTable: table, x: 7, y: 0, isBlack: true}
    table[7][7] = &Rook{cTable: table, x: 7, y: 7, isBlack: true}

    // Knight
    table[0][1] = &Knight{cTable: table, x: 0, y: 1, isBlack: false}
    table[0][6] = &Knight{cTable: table, x: 0, y: 6, isBlack: false}
    table[7][1] = &Knight{cTable: table, x: 7, y: 1, isBlack: true}
    table[7][6] = &Knight{cTable: table, x: 7, y: 6, isBlack: true}

    // Bishop
    table[0][2] = &Bishop{cTable: table, x: 0, y: 2, isBlack: false}
    table[0][5] = &Bishop{cTable: table, x: 0, y: 5, isBlack: false}
    table[7][2] = &Bishop{cTable: table, x: 7, y: 2, isBlack: true}
    table[7][5] = &Bishop{cTable: table, x: 7, y: 5, isBlack: true}

    // Queen
    table[0][3] = &Queen{cTable: table, x: 0, y: 3, isBlack: false}
    table[7][3] = &Queen{cTable: table, x: 7, y: 3, isBlack: true}

    // King
    table[0][4] = &King{cTable: table, x: 0, y: 4, isBlack: false}
    table[7][4] = &King{cTable: table, x: 7, y: 4, isBlack: true}

    return table
}

// This functions checks black and white winner states
func checkBlackWin(table *ChessTable) bool {
    // Find king
    var king Figure
    var kingx int
    var kingy int

    for i := 0; i <= 7; i++ {
        for j := 0; j <= 7; j++ {
            if k, ok := table[i][j].(*King); ok && !k.IsFigureBlack() {
                king = k
                kingx = k.x
                kingy = k.y
            }
        }
    }

    if (king == nil) { return true; }
    whereKingMaybeOffsets := make(map[[2]int]bool)
    whereKingMaybeOffsets[[2]int{0, 0}] = false

    for i := -1; i <= 1; i++ {
        for j := -1; j <= 1; j++ {
            if i == 0 && j == 0 { continue; }
            if king.CanFigurePass(i, j) {
                whereKingMaybeOffsets[[2]int{i, j}] = false
            }
        }
    }

    for offset, _ := range whereKingMaybeOffsets {
        // Check if every black figure can beat king
        for i := 0; i <= 7; i++ {
            for j := 0; j <= 7; j++ {
                if table[i][j] != nil && table[i][j].IsFigureBlack() {
                    canBeat := table[i][j].CanFigurePass(i - (kingx - offset[0]), j - (kingy - offset[1]))
                    if canBeat { whereKingMaybeOffsets[offset] = true }
                }
            }
        }
    }

    // If all whereKingMaybeOffsets element are true, so
    // the king has nowhere to go
    for _, value := range whereKingMaybeOffsets {
        if !value {
            return false
        }
    }
    return true
}

func checkWhiteWin(table *ChessTable) bool {
    // Find king
    var king Figure
    var kingx int
    var kingy int
    for i := 7; i >= 0; i-- {
        for j := 7; j >= 0; j-- {
            if k, ok := table[i][j].(*King); ok && k.IsFigureBlack() {
                king = k
                kingx = k.x
                kingy = k.y
            }
        }
    }

    if (king == nil) { return true; }
    whereKingMaybeOffsets := make(map[[2]int]bool)
    whereKingMaybeOffsets[[2]int{0, 0}] = false

    for i := -1; i <= 1; i++ {
        for j := -1; j <= 1; j++ {
            if i == 0 && j == 0 { continue; }
            if king.CanFigurePass(i, j) {
                whereKingMaybeOffsets[[2]int{i, j}] = false
            }
        }
    }

    for offset, _ := range whereKingMaybeOffsets {
        // Check if every black figure can beat king
        for i := 0; i <= 7; i++ {
            for j := 0; j <= 7; j++ {
                if table[i][j] != nil && !table[i][j].IsFigureBlack() {
                    canBeat := table[i][j].CanFigurePass((kingx + offset[0]) - i, (kingy + offset[1]) - j)
                    if canBeat { whereKingMaybeOffsets[offset] = true }
                }
            }
        }
    }

    // If all whereKingMaybeOffsets element are true, so
    // the king has nowhere to go
    for _, value := range whereKingMaybeOffsets {
        if !value {
            return false
        }
    }
    return true
}

// Function to save game data to database
func saveResultsToDatabase(session *GameSession, players []string, isDraw bool, winnerLogin string, timeStarted time.Time) {
    if session.ExtIdentifier == "" {
        gameHistoryObject := database.GamesHistory{
            GameStartTimestamp: timeStarted,
            GameEndTimestamp: time.Now(),
            IsDraw: isDraw,
            WinnerLoginKey: winnerLogin,
        }

        // Ignore errors
        database.DB.Save(&gameHistoryObject)

        // Add new player list objects
        for _, v := range players {
            pObject := database.PlayerList{
                LoginKey: v,
                GamesHistoryKey: gameHistoryObject.ID,
            }

            database.DB.Save(&pObject)
        }
    }
}

// Function to close all connections
func closeConnections(session *GameSession) {
    for _, v := range session.GetAllPlayers() {
        v.Connection.CloseConnection("Connection closed")
        session.deleteConnection(v)
    }
}

// Function to broadcast json messages
func broadcastJSONMessages(message interface{}, currConnection *GameClientConnection, session *GameSession) {
    data, err := json.Marshal(message)
    if err != nil {
        panic(err)
    }

    for _, v := range session.GetAllPlayers() {
        if v != currConnection {
            v.Connection.GetConnection().WriteMessage(websocket.TextMessage, data)
        }
    }
}

