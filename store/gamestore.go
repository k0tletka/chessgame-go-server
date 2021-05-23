package store

import (
    "sync"
    "errors"
)

var (
    // Errors
    GameNotFoundError = errors.New("Game with the specified id is not found")

    // Game stores
    GameStores = []GameStore{}

    // Utility
    idGameCounter = 0
    gameStoresMutex = sync.RWMutex{}
)

// Type to input turns
type GameTurn struct {
    FigposX int
    FigposY int
    AltX int
    AltY int
    Surrender bool
}

// This type represents a game
type GameStore struct {
    GameID int
    PlayerOneLogin string
    PlayerTwoLogin string
    GameStarted bool
    IsPlayerOneTurn bool
    GameTitle string
    AckChannel chan string
    SendTurnRequest chan *GameTurn
    SendTurnResponse chan string
}

// Functions to add, remove and getting game by id
func RegisterNewGameStore(gametitle string, player string) int {

    // Thread-safe operation
    gameStoresMutex.Lock()
    defer gameStoresMutex.Unlock()

    idGameCounter++

    GameStores = append(GameStores, GameStore{
        GameID: idGameCounter,
        PlayerOneLogin: player,
        PlayerTwoLogin: "",
        GameStarted: false,
        IsPlayerOneTurn: true,
        GameTitle: gametitle,
        AckChannel: make(chan string),
        SendTurnRequest: make(chan *GameTurn, 1),
        SendTurnResponse: make(chan string, 1),
    })

    return idGameCounter
}

func RemoveGameStore(id int) error {

    // Thread-safe operation
    gameStoresMutex.Lock()
    defer gameStoresMutex.Unlock()

    for i, proc := range GameStores {
        if proc.GameID == id {
            close(GameStores[i].AckChannel)
            close(GameStores[i].SendTurnRequest)
            close(GameStores[i].SendTurnResponse)

            GameStores[i] = GameStores[len(GameStores) - 1]
            GameStores = GameStores[:len(GameStores) - 1]

            return nil
        }
    }

    return GameNotFoundError
}

func GetGameStore(id int) (*GameStore, error) {

    // Thread-safe operation
    gameStoresMutex.RLock()
    defer gameStoresMutex.RUnlock()

    for i, proc := range GameStores {
        if proc.GameID == id {
            return &GameStores[i], nil
        }
    }

    return nil, GameNotFoundError
}
