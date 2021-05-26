package game

import (
    "encoding/json"
    "errors"
    "sync"

    ws "GoChessgameServer/websocket"

    "github.com/gorilla/websocket"
)

var (
    // Errors
    GameNotFoundError = errors.New("Game with the specified id is not found")
    UserAlreadyConnected = errors.New("User is already connected to this game")
    NotEnoughtPlayers = errors.New("Not enought players for session, min 2 required")
    SessionFullError = errors.New("Can't connect to session: session is full")

    // Game stores
    SessionStore = CreateNewSessionList()
)

// This type represent active game connection of client
type GameClientConnection struct {
    Connection      *ws.WebsocketConnection
    Login           string

    // This flags means that connected user is from another instance
    ExternalUser    bool

    // Channel for transmitting websocket requests to game
    ReadChannel     chan *Turn
}

// Type represents linked list for clients turn loop
type GameTurner struct {
    GameConnection  *GameClientConnection
    Next            *GameTurner
}

// This type stores information about game
// and performs connectivity between clients
type GameSession struct {
    // Local instance game information
    GameID              int
    GameStarted         bool
    GameTitle           string
    PlayersMax          int

    // External instance identifier
    ExtIdentifier       string

    PlayerTurner        *GameTurner
    playersMutex        *sync.RWMutex
    players             []*GameClientConnection
}

func (g *GameSession) AddNewConnection(conn *GameClientConnection) error {
    g.playersMutex.Lock()
    defer g.playersMutex.Unlock()

    // Check if connection already exists
    for _, v := range g.players {
        if v.Connection == conn.Connection {
            return UserAlreadyConnected
        }
    }

    if g.PlayersMax - len(g.players) >= 1 {
        g.players = append(g.players, conn)
    } else {
        return SessionFullError
    }

    // Add close connection handler to delete connection from list
    g.setCloseHandler(conn)
    g.addConnectionToTurner(g.PlayerTurner, conn)

    if g.PlayersMax - len(g.players) == 0 {
        go g.StartSession()
    }

    return nil
}

// Get all player of session
func (g *GameSession) GetAllPlayers() []*GameClientConnection {
    return g.players
}

// Method for deleting connection from list. Used in connection close handler
func (g *GameSession) deleteConnection(conn *GameClientConnection) {
    g.playersMutex.Lock()
    defer g.playersMutex.Unlock()

    for i, v := range g.players {
        if v == conn {
            g.players[i] = g.players[len(g.players) - 1]
            g.players = g.players[:len(g.players) - 1]

            g.deleteConnectionFromTurner(nil, g.PlayerTurner, conn)

            close(conn.ReadChannel)

            return
        }
    }
}

func (g *GameSession) addConnectionToTurner(turner *GameTurner, conn *GameClientConnection) {
    if turner == nil {
        g.PlayerTurner = &GameTurner{GameConnection: conn}
        return
    }

    if turner.Next != nil {
        g.addConnectionToTurner(turner.Next, conn)
    } else {
        turner.Next = &GameTurner{GameConnection: conn}
    }
}

func (g *GameSession) deleteConnectionFromTurner(parent *GameTurner, turner *GameTurner, conn *GameClientConnection) {
    if turner == nil {
        return
    }

    if turner.GameConnection == conn {
        if parent == nil { g.PlayerTurner = turner.Next; }
        parent.Next = turner.Next
    } else {
        g.deleteConnectionFromTurner(turner, turner.Next, conn)
    }
}

func (g *GameSession) setCloseHandler(conn *GameClientConnection) {
    conn.Connection.AddCloseHandler(func (wc *ws.WebsocketConnection) {
        g.deleteConnection(conn)

        g.playersMutex.Lock()
        defer g.playersMutex.Unlock()

        // If it was last connection, delete session from list
        if len(g.players) == 0 {
            SessionStore.RemoveGameSession(g.GameID)
        }
    })
}


func (g *GameSession) StartSession() {
    // Remind all players that game has started
    request := struct{
        GameStarted bool `json:"game_started"`
    }{
        GameStarted: true,
    }

    if data, err := json.Marshal(&request); err != nil {
        panic(err)
    } else {
        for _, c := range g.players {
            c.Connection.GetConnection().WriteMessage(websocket.TextMessage, data)
        }

        controlGame(g)
    }
}

// This type represents game session list
type GameSessionList struct {
    sync.RWMutex

    idGameCounter   int
    sessions        []*GameSession
}

func CreateNewSessionList() *GameSessionList {
    return &GameSessionList{
        sessions: []*GameSession{},
    }
}

// Function to register new game session
func (g *GameSessionList) RegisterNewGameSession(gameTitle string, userInfo *GameClientConnection, playersMax int) (int, error) {
    if playersMax < 2 {
        return 0, NotEnoughtPlayers
    }

    g.Lock()

    g.idGameCounter++
    newSession := &GameSession{
        GameID: g.idGameCounter,
        GameStarted: false,
        GameTitle: gameTitle,
        PlayersMax: playersMax,
        // TODO: Make registering external identifier for DHT-network
        PlayerTurner: nil,
        playersMutex: &sync.RWMutex{},
        players: []*GameClientConnection{},
    }

    newSession.AddNewConnection(userInfo)
    g.sessions = append(g.sessions, newSession)

    g.Unlock()

    return g.idGameCounter, nil
}

// Function to get appropriate game session by id
func (g *GameSessionList) GetGameSession(gameId int) (*GameSession, error) {
    g.RLock()
    defer g.RUnlock()

    for _, v := range g.sessions {
        if v.GameID == gameId {
            return v, nil
        }
    }

    return nil, GameNotFoundError
}

func (g *GameSessionList) RemoveGameSession(gameId int) {
    g.Lock()
    defer g.Unlock()

    for i, v := range g.sessions {
        if v.GameID == gameId {
            g.sessions[i] = g.sessions[len(g.sessions) - 1]
            g.sessions = g.sessions[:len(g.sessions) - 1]

            return
        }
    }
}

func (g *GameSessionList) GetAllGameSessions() []*GameSession {
    return g.sessions
}
