package store

import (
    "crypto/rand"
    "log"
    "os"
    "sync"
    "errors"

    "GoChessgameServer/logger"

    jwt "github.com/dgrijalva/jwt-go"
    lp "github.com/jcuga/golongpoll"
)

// Store logger
var storeLogger *log.Logger

// This string contains a key, stat will be
// used for jwt signing
var JWTKey []byte

var idGameCounter = 0
var GameStores = []GameStore{}

// Mutexs for thread-safe operations on
// game processes
var createGameMutex = sync.Mutex{}
var removeGameMutex = sync.Mutex{}

// Errors
var GameNotFoundError = errors.New("Game with the specified id is not found")

// This type represent token claims with login string
type JWTClaims struct {
    Login string
    IsAdmin bool
    jwt.StandardClaims
}

// This type represents a game
type GameStore struct {
    GameID int
    PlayerOneLogin string
    PlayerTwoLogin string
    EndGameManager *lp.LongpollManager
    SendTurnManager *lp.LongpollManager
}

func init() {
    // Create store logger
    storeLogger = logger.AddNewLogger("Store", os.Stdout, log.LstdFlags | log.Lmsgprefix)

    // Generate JWT key
    JWTKey = make([]byte, 256, 256)
    _, err := rand.Read(JWTKey)
    if err != nil {
        storeLogger.Fatalf("Error when generating key: %s", err.Error())
    }
}

// Functions to add, remove and getting game by id
func RegisterNewGameStore(playerone string, playertwo string) int {

    // Thread-safe operation
    createGameMutex.Lock()

    endGameManager, _ := lp.StartLongpoll(lp.Options{})
    sendTurnManager, _ := lp.StartLongpoll(lp.Options{})
    idGameCounter++

    _ = append(GameStores, GameStore{
        GameID: idGameCounter,
        PlayerOneLogin: playerone,
        PlayerTwoLogin: playertwo,
        EndGameManager: endGameManager,
        SendTurnManager: sendTurnManager,
    })

    createGameMutex.Unlock()
    return idGameCounter
}

func RemoveGameStore(id int) error {

    // Thread-safe operation
    removeGameMutex.Lock()

    for i, proc := range GameStores {
        if proc.GameID == id {
            GameStores[i] = GameStores[len(GameStores) - 1]
            GameStores = GameStores[:len(GameStores) - 1]

            removeGameMutex.Unlock()
            return nil
        }
    }

    removeGameMutex.Unlock()
    return GameNotFoundError
}

func GetGameStore(id int) (*GameStore, error) {

    for _, proc := range GameStores {
        if proc.GameID == id {
            return &proc, nil
        }
    }

    return nil, GameNotFoundError
}
