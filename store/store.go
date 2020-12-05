package store

import (
    "crypto/rand"
    "log"
    "os"
    "sync"
    "errors"
    "io/ioutil"

    "GoChessgameServer/logger"
    c "GoChessgameServer/conf"

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
var gameStoresMutex = sync.Mutex{}

// Longpolls
var WaitGameLM *lp.LongpollManager
var WaitTurnLM *lp.LongpollManager
var EndGameLM *lp.LongpollManager

// String variable for markdown motd storing
var MotdString string

// Errors
var GameNotFoundError = errors.New("Game with the specified id is not found")

// This type represent token claims with login string
type JWTClaims struct {
    Login string
    IsAdmin bool
    jwt.StandardClaims
}

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

func init() {
    // Create store logger
    storeLogger = logger.AddNewLogger("Store", os.Stdout, log.LstdFlags | log.Lmsgprefix)

    // Generate JWT key
    JWTKey = make([]byte, 256, 256)
    _, err := rand.Read(JWTKey)
    if err != nil {
        storeLogger.Fatalf("Error when generating key: %s", err.Error())
    }

    // Load markdown file
    mdFile := c.Conf.App.MarkdownFile
    if mdFile == "" {
        storeLogger.Println("Markdown file is not set, skipping...")
        MotdString = "Motd file is not set on the server."
        return
    }

    mdFd, err := os.OpenFile(mdFile, os.O_RDONLY | os.O_CREATE, 0755)
    if err != nil {
        storeLogger.Fatalln(err)
    }
    defer mdFd.Close()

    // Read markdown into string variable
    readedMotd, err := ioutil.ReadAll(mdFd)
    if err != nil {
        storeLogger.Fatalln(err)
    }
    MotdString = string(readedMotd)

    // Initialize longpoll managers
    WaitGameLM, err = lp.StartLongpoll(lp.Options{
        EventTimeToLiveSeconds: 3,
        DeleteEventAfterFirstRetrieval: true,
        MaxLongpollTimeoutSeconds: 3600,
    })
    if err != nil {
        storeLogger.Fatalln(err)
    }

    WaitTurnLM, err = lp.StartLongpoll(lp.Options{
        EventTimeToLiveSeconds: 3,
        DeleteEventAfterFirstRetrieval: true,
        MaxLongpollTimeoutSeconds: 600,
    })
    if err != nil {
        storeLogger.Fatalln(err)
    }

    EndGameLM, err = lp.StartLongpoll(lp.Options{
        EventTimeToLiveSeconds: 3600 * 24,
        DeleteEventAfterFirstRetrieval: false,
        MaxLongpollTimeoutSeconds: 3600 * 24,
    })
    if err != nil {
        storeLogger.Fatalln(err)
    }
}

// Functions to add, remove and getting game by id
func RegisterNewGameStore(gametitle string, player string) int {

    // Thread-safe operation
    gameStoresMutex.Lock()
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

    gameStoresMutex.Unlock()
    return idGameCounter
}

func RemoveGameStore(id int) error {

    // Thread-safe operation
    gameStoresMutex.Lock()

    for i, proc := range GameStores {
        if proc.GameID == id {
            close(GameStores[i].AckChannel)
            close(GameStores[i].SendTurnRequest)
            close(GameStores[i].SendTurnResponse)

            GameStores[i] = GameStores[len(GameStores) - 1]
            GameStores = GameStores[:len(GameStores) - 1]

            gameStoresMutex.Unlock()
            return nil
        }
    }

    gameStoresMutex.Unlock()
    return GameNotFoundError
}

func GetGameStore(id int) (*GameStore, error) {

    for i, proc := range GameStores {
        if proc.GameID == id {
            return &GameStores[i], nil
        }
    }

    return nil, GameNotFoundError
}
