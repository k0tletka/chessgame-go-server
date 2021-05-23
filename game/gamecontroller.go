package game

import (
    //"encoding/json"
    "time"
    //"strconv"

    "GoChessgameServer/store"
    "GoChessgameServer/database"
)

const queryResString = `INSERT INTO
    dbo.GamesHistory(GameStartTimestamp, GameEndTimestamp, IsDraw, WinnerLogin, PlayerOneLogin, PlayerTwoLogin)
VALUES
    ($1, $2, $3, $4, $5, $6)`

type notifyWaitGameData struct {
    TurnMade bool `json:"turnMade"`
    FigposX int `json:"figposx"`
    FigposY int `json:"figposy"`
    AltX int `json:"xalt"`
    AltY int `json:"yalt"`
    Surrendered bool `json:"surrendered"`
}

// This function is a main entrypoint handler,
// that controls game process.
// This function starts as coroutine
func ControlGame(gameStore *store.GameStore) {

    res := struct{
        WinnerLogin string `json:"winned"`
        Draw bool `json:"isdraw"`
    }{}
    table := ChessTable{}
    gameStarted := time.Now()

    // Create game session
    gameSession := GameSession{
        gameStore: gameStore,
        cTable: &table,
    }

    // Create instance of game result for database
    databaseResult := database.GamesHistory{
        GameStartTimestamp: time.Now(),
        PlayerOneLoginKey: gameStore.PlayerOneLogin,
        PlayerTwoLoginKey: gameStore.PlayerTwoLogin,
    }

    // Fill table with figures
    fillTable(gameSession)

    // Handle turns
    turnTimer := time.NewTimer(time.Minute * 5)
    for {
        select {
        case <-turnTimer.C:
            // Timer out, it is draw
            res.Draw = true

            //jsonBytes, _ := json.Marshal(&res)
            //store.EndGameLM.Publish(strconv.Itoa(gameStore.GameID), string(jsonBytes))

            //notifyData := notifyWaitGameData{TurnMade: false}
            //notifyPlayerWait(gameStore.GameID, &notifyData)
            _ = store.RemoveGameStore(gameStore.GameID)

            // Execute database query
            saveResultToDatabase(&databaseResult, true, "", gameStarted)
            return
        case turn := <-gameStore.SendTurnRequest:
            // Handle turn
            if turn.Surrender {
                // Player surrendered, write results
                if gameStore.IsPlayerOneTurn {
                    res.WinnerLogin = gameStore.PlayerTwoLogin
                } else {
                    res.WinnerLogin = gameStore.PlayerOneLogin
                }

                //jsonBytes, _ := json.Marshal(&res)
                //store.EndGameLM.Publish(strconv.Itoa(gameStore.GameID), string(jsonBytes))

                gameStore.SendTurnResponse <- ""
                //notifyData := notifyWaitGameData{TurnMade: true, Surrendered: true}
                //notifyPlayerWait(gameStore.GameID, &notifyData)
                _ = store.RemoveGameStore(gameStore.GameID)

                // Execute database query
                saveResultToDatabase(&databaseResult, false, res.WinnerLogin, gameStarted)
                return
            }

            // Check figure existent
            if table[turn.FigposX][turn.FigposY] == nil {
                gameStore.SendTurnResponse <- "There is no figure"
                continue
            }

            // Check for player turn and perform operation
            if (gameStore.IsPlayerOneTurn && table[turn.FigposX][turn.FigposY].IsFigureBlack()) || (!gameStore.IsPlayerOneTurn && !table[turn.FigposX][turn.FigposY].IsFigureBlack()) {
                gameStore.SendTurnResponse <- "You can't move figure that you don't own"
                continue
            }

            if !table[turn.FigposX][turn.FigposY].CanFigurePass(turn.AltX, turn.AltY) {
                gameStore.SendTurnResponse <- "You can't move this figure at the specified location"
                continue
            }


            // Pass table[turn.FigposX][turn.FigposY]
            _ = table[turn.FigposX][turn.FigposY].Pass(turn.AltX, turn.AltY)

            // Check win state
            if checkBlackWin(gameSession) {
                // Prepare for sending response on channels
                res.WinnerLogin = gameStore.PlayerTwoLogin
                //jsonBytes, _ := json.Marshal(&res)
                //store.EndGameLM.Publish(strconv.Itoa(gameStore.GameID), string(jsonBytes))

                // Send json to end game longpoll
                gameStore.SendTurnResponse <- ""

                //notifyData := notifyWaitGameData{TurnMade: true, FigposX: turn.FigposX, FigposY: turn.FigposY, AltX: turn.AltX, AltY: turn.AltY}
                //notifyPlayerWait(gameStore.GameID, &notifyData)
                _ = store.RemoveGameStore(gameStore.GameID)

                // Write results to database
                saveResultToDatabase(&databaseResult, false, res.WinnerLogin, gameStarted)
                return
            }

            if checkWhiteWin(gameSession) {
                // Prepare for sending response on channels
                res.WinnerLogin = gameStore.PlayerOneLogin
                //jsonBytes, _ := json.Marshal(&res)
                //store.EndGameLM.Publish(strconv.Itoa(gameStore.GameID), string(jsonBytes))

                // Send json to end game longpoll
                gameStore.SendTurnResponse <- ""

                //notifyData := notifyWaitGameData{TurnMade: true, FigposX: turn.FigposX, FigposY: turn.FigposY, AltX: turn.AltX, AltY: turn.AltY}
                //notifyPlayerWait(gameStore.GameID, &notifyData)
                _ = store.RemoveGameStore(gameStore.GameID)

                // Write results to database
                saveResultToDatabase(&databaseResult, false, res.WinnerLogin, gameStarted)
                return
            }

            // Turn made, send notification to opponent and wait ack
            //notifyData := notifyWaitGameData{TurnMade: true, FigposX: turn.FigposX, FigposY: turn.FigposY, AltX: turn.AltX, AltY: turn.AltY}
            //notifyPlayerWait(gameStore.GameID, &notifyData)
            ackTimer := time.NewTimer(time.Second * 3)

            select {
            case <-ackTimer.C:
                // Opponent not response, aborting game
                if gameStore.IsPlayerOneTurn {
                    res.WinnerLogin = gameStore.PlayerOneLogin
                } else {
                    res.WinnerLogin = gameStore.PlayerTwoLogin
                }

                //jsonBytes, _ := json.Marshal(&res)
                //store.EndGameLM.Publish(strconv.Itoa(gameStore.GameID), string(jsonBytes))
                gameStore.SendTurnResponse <- ""
                _ = store.RemoveGameStore(gameStore.GameID)

                // Write results to database
                saveResultToDatabase(&databaseResult, false, res.WinnerLogin, gameStarted)
                return
            case <-gameStore.AckChannel:
                turnTimer = time.NewTimer(time.Minute * 5)
                if gameStore.IsPlayerOneTurn {
                    gameStore.IsPlayerOneTurn = false
                } else {
                    gameStore.IsPlayerOneTurn = true
                }

                gameStore.SendTurnResponse <- ""
            }
        }
    }
}

// This function fills table with figures
func fillTable(session GameSession) {
    // Pawn
    for i := 0; i <= 7; i++ {
        session.cTable[1][i] = &Pawn{cTable: session.cTable, x: 1, y: i, isBlack: false, isMoved: false}
        session.cTable[6][i] = &Pawn{cTable: session.cTable, x: 6, y: i, isBlack: true, isMoved: false}
    }

    // Rook
    session.cTable[0][0] = &Rook{cTable: session.cTable, x: 0, y: 0, isBlack: false}
    session.cTable[0][7] = &Rook{cTable: session.cTable, x: 0, y: 7, isBlack: false}
    session.cTable[7][0] = &Rook{cTable: session.cTable, x: 7, y: 0, isBlack: true}
    session.cTable[7][7] = &Rook{cTable: session.cTable, x: 7, y: 7, isBlack: true}

    // Knight
    session.cTable[0][1] = &Knight{cTable: session.cTable, x: 0, y: 1, isBlack: false}
    session.cTable[0][6] = &Knight{cTable: session.cTable, x: 0, y: 6, isBlack: false}
    session.cTable[7][1] = &Knight{cTable: session.cTable, x: 7, y: 1, isBlack: true}
    session.cTable[7][6] = &Knight{cTable: session.cTable, x: 7, y: 6, isBlack: true}

    // Bishop
    session.cTable[0][2] = &Bishop{cTable: session.cTable, x: 0, y: 2, isBlack: false}
    session.cTable[0][5] = &Bishop{cTable: session.cTable, x: 0, y: 5, isBlack: false}
    session.cTable[7][2] = &Bishop{cTable: session.cTable, x: 7, y: 2, isBlack: true}
    session.cTable[7][5] = &Bishop{cTable: session.cTable, x: 7, y: 5, isBlack: true}

    // Queen
    session.cTable[0][3] = &Queen{cTable: session.cTable, x: 0, y: 3, isBlack: false}
    session.cTable[7][3] = &Queen{cTable: session.cTable, x: 7, y: 3, isBlack: true}

    // King
    session.cTable[0][4] = &King{cTable: session.cTable, x: 0, y: 4, isBlack: false}
    session.cTable[7][4] = &King{cTable: session.cTable, x: 7, y: 4, isBlack: true}
}

// This functions checks black and white winner states
func checkBlackWin(session GameSession) bool {
    // Find king
    var king Figure
    var kingx int
    var kingy int

    for i := 0; i <= 7; i++ {
        for j := 0; j <= 7; j++ {
            if k, ok := session.cTable[i][j].(*King); ok && !k.IsFigureBlack() {
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
                if session.cTable[i][j] != nil && session.cTable[i][j].IsFigureBlack() {
                    canBeat := session.cTable[i][j].CanFigurePass(i - (kingx - offset[0]), j - (kingy - offset[1]))
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

func checkWhiteWin(session GameSession) bool {
    // Find king
    var king Figure
    var kingx int
    var kingy int
    for i := 7; i >= 0; i-- {
        for j := 7; j >= 0; j-- {
            if k, ok := session.cTable[i][j].(*King); ok && k.IsFigureBlack() {
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
                if session.cTable[i][j] != nil && !session.cTable[i][j].IsFigureBlack() {
                    canBeat := session.cTable[i][j].CanFigurePass((kingx + offset[0]) - i, (kingy + offset[1]) - j)
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

// This function sends to waiting player a notification
// that opponent's turn has been ended
//func notifyPlayerWait(gameID int, resp *notifyWaitGameData) {
    //jsonBytes, _ := json.Marshal(&resp)
    //store.WaitTurnLM.Publish(strconv.Itoa(gameID), string(jsonBytes))
//}

// Function to save game data to database
func saveResultToDatabase(resultObject *database.GamesHistory, isDraw bool, winnerLogin string, timeEnded time.Time) {
    resultObject.IsDraw = isDraw
    resultObject.WinnerLoginKey = winnerLogin
    resultObject.GameEndTimestamp = timeEnded

    // Ignore errors
    database.DB.Save(resultObject)
}
