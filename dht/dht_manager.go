package dht

import (
    "net/http"
    "net"
    "fmt"
    "context"
    "time"
    "encoding/json"
    "encoding/hex"
    "errors"

    ws "GoChessgameServer/websocket"
    c "GoChessgameServer/conf"
    u "GoChessgameServer/util"
    "GoChessgameServer/database"
    "GoChessgameServer/auth"

    "github.com/gorilla/websocket"
)

var (
    // Errors
    HostNotFound = errors.New("Host requested by given identifier not found")
    TimeoutError = errors.New("Connection has timed out")
)

// This type represent static peer from configuration
type StaticPeerConnection struct {
    Connection *ws.WebsocketConnection
    StaticPeer *c.SPeer
}

// This is main manager of server DHT API,
// that performs handshaking process and provides
// programm API for other modules
type DHTManager struct {
    // Store of server connections
    wsConns *ws.WebsocketStore

    // List with connect information about static peers and connection object
    staticPeerConnections []*StaticPeerConnection

    // List of connections, those data stored in database
    databasePeerConnections map[string]*ws.WebsocketConnection

    // Read synchronization channel to wait for connection response
    readSynchronizationChannels map[*ws.WebsocketConnection]chan *dhtAPIBaseRequest

    // Websocket dialer
    wsDialer *websocket.Dialer
}

func CreateNewDHTManager() *DHTManager {
    wsDialer := &websocket.Dialer{
        HandshakeTimeout: 3 * time.Second,
    }

    newManager := &DHTManager {
        wsConns: &ws.WebsocketStore{},
        staticPeerConnections: []*StaticPeerConnection{},
        wsDialer: wsDialer,
    }

    // Fill static peer connections
    newManager.fillStaticPeerConnections()

    // Start handshake goroutine
    go newManager.startHandshakeProcedure()

    return newManager
}

// Public program API
func (m *DHTManager) GetServerIdentifier() string {
    return hex.EncodeToString(dhtServerIdentifier[:])
}

// Function to get information about host by its server identifier
func (m *DHTManager) GetHostInfoByServerIdentifier(serverIdentifier string) (*DHTHostInformation, error) {
    if v, ok := m.databasePeerConnections[serverIdentifier]; ok && v != nil && !v.Closed() {
        // Send request to get information about host
        baseRequest := dhtAPIBaseRequest{
            MethodName: "hostinfo",
        }

        var data []byte
        var err error

        if data, err = json.Marshal(&baseRequest); err != nil {
            return nil, err
        }

        // Init channel for result awaiting
        readChan := make(chan *dhtAPIBaseRequest, 1)
        m.readSynchronizationChannels[v] = readChan

        v.GetConnection().WriteMessage(websocket.TextMessage, data)

        // Wait for result with timeout
        timer := time.NewTimer(5 * time.Second)

        select {
        case response := <-readChan:
            delete(m.readSynchronizationChannels, v)
            close(readChan)

            result := &DHTHostInformation{}

            // Parse response
            if err = json.Unmarshal(response.Args, result); err != nil {
                return nil, err
            }

            result.ClientAPIIPAddress = v.GetConnection().RemoteAddr().(*net.TCPAddr).IP.String()
            result.GameAPIIPAddress = v.GetConnection().RemoteAddr().(*net.TCPAddr).IP.String()

            return result, nil
        case <-timer.C:
            delete(m.readSynchronizationChannels, v)
            close(readChan)

            return nil, TimeoutError
        }

    }

    return nil, HostNotFound
}

// Method for verifying server token
func (m *DHTManager) VerifyServerToken(tokenString, serverIdentifier string) (*auth.JWTUserClaim, bool, error) {
    if v, ok := m.databasePeerConnections[serverIdentifier]; ok && v != nil && !v.Closed() {
        // Send request to verify token
        var data []byte
        var err error

        request := struct{
            TokenToVerify   string  `json:"token"`
        }{
            TokenToVerify: tokenString,
        }

        if data, err = json.Marshal(&request); err != nil {
            return nil, false, err
        }

        baseRequest := dhtAPIBaseRequest{
            MethodName: "verifytoken",
            Args: data,
        }

        if data, err = json.Marshal(&baseRequest); err != nil {
            return nil, false, err
        }

        // Init channel for result awaiting
        readChan := make(chan *dhtAPIBaseRequest, 1)
        m.readSynchronizationChannels[v] = readChan

        v.GetConnection().WriteMessage(websocket.TextMessage, data)

        // Wait for result with timeout
        timer := time.NewTimer(5 * time.Second)

        select {
        case response := <-readChan:
            delete(m.readSynchronizationChannels, v)
            close(readChan)

            result := struct{
                Verified    bool            `json:"verified"`
                TokenData   json.RawMessage `json:"token_data"`
            }{}

            // Parse response
            if err = json.Unmarshal(response.Args, &result); err != nil {
                return nil, false, err
            }

            tokenData := &auth.JWTUserClaim{}

            if err = json.Unmarshal(result.TokenData, &tokenData); err != nil {
                return nil, false, err
            }

            return tokenData, result.Verified, nil
        case <-timer.C:
            delete(m.readSynchronizationChannels, v)
            close(readChan)

            return nil, false, TimeoutError
        }
    }

    return nil, false, HostNotFound
}

// Private methods
func (m *DHTManager) fillStaticPeerConnections() {
    for _, v := range c.Conf.StaticPeers {
        m.staticPeerConnections = append(m.staticPeerConnections, &StaticPeerConnection{
            Connection: nil,
            StaticPeer: &v,
        })
    }
}

func (m *DHTManager) startHandshakeProcedure() {
    var handshakeTimeout uint

    if !c.DecodeMetadata.IsDefined("dht_api", "handshake_period") {
        handshakeTimeout = 60 // 1 minute for default
    } else {
        handshakeTimeout = c.Conf.DHTApi.HandshakePeriod
    }

    var connectionLimit uint

    if !c.DecodeMetadata.IsDefined("dht_api", "connections_limit") {
        connectionLimit = 5 // 5 is default
    } else {
        connectionLimit = c.Conf.DHTApi.ConnectionsLimit
    }

    for {
        // Enstablish connection with static peers and send them request with out
        // server identifier
        for _, v := range m.staticPeerConnections {
            var conn *ws.WebsocketConnection

            // If connection is nil or connection closed, create new connection
            if v.Connection == nil || v.Connection.Closed() {
                wsConn, err := m.createConnection(
                    v.StaticPeer.ServerName,
                    v.StaticPeer.ConnectionPort,
                    v.StaticPeer.UseTLS,
                )

                if err != nil {
                    dhtLogger.Printf("Connection to %s peer failed: %s\n", v.StaticPeer.ServerName, err.Error())
                    continue
                }

                // Add connection to list
                conn = ws.NewWebsocketConnection(wsConn, m.connectionReadHandler, m.wsConns)
                v.Connection = conn
            } else {
                conn = v.Connection
            }

            m.sendHandshakeRequest(conn, connectionLimit)
        }

        // Enstablish connections with database hosts
        databaseHosts := []database.DHTHosts{}

        database.DB.
            Where("is_peer_static = ?", false).
            Where("srv_local_identifier = ?", dhtServerIdentifier[:]).
            Find(&databaseHosts)

        for _, v := range databaseHosts {
            encodedIdentifier := hex.EncodeToString(v.ServerIdentifier)

            // If host connected already, just use connection
            var conn *ws.WebsocketConnection

            if mConn, ok := m.databasePeerConnections[encodedIdentifier]; ok && mConn != nil && !mConn.Closed() {
                conn = mConn
            } else {
                wsConn, err := m.createConnection(
                    v.IPAddress,
                    v.Port,
                    v.UseTLS,
                )

                if err != nil {
                    dhtLogger.Printf("Connection to %s peer failed: %s\n", v.IPAddress, err.Error())
                    continue
                }

                conn = ws.NewWebsocketConnection(wsConn, m.connectionReadHandler, m.wsConns)
                m.databasePeerConnections[encodedIdentifier] = conn
            }

            m.sendHandshakeRequest(conn, connectionLimit)
        }

        // Sleep until next handshake timeout
        time.Sleep(time.Duration(handshakeTimeout) * time.Second)
    }
}

func (m *DHTManager) sendHandshakeRequest(conn *ws.WebsocketConnection, connectionLimit uint) {
    // Get listening port of server api
    _, listenport := u.GetListenInformationServerAPI()

    // Connection exist, so just send handshake message
    request := struct{
        ServerIdentifier    string  `json:"server_identifier"`
        ServerAPIPort       uint16  `json:"server_api_port"`
        UseTLS              bool    `json:"use_tls"`
        ConnectionLimit     uint    `json:"connection_limit"`
    }{
        ServerIdentifier: m.GetServerIdentifier(),
        ServerAPIPort: listenport,
        UseTLS: c.Conf.DHTApi.UseTLS,
        ConnectionLimit: connectionLimit,
    }

    var data []byte
    var err error

    if data, err = json.Marshal(&request); err != nil {
        dhtLogger.Printf("Error when marshalling request to peer %s: %s\n", conn.GetConnection().RemoteAddr(), err.Error())
        return
    }

    baseRequest := dhtAPIBaseRequest{
        MethodName: "handshake",
        Args: data,
    }

    if data, err = json.Marshal(&baseRequest); err != nil {
        dhtLogger.Printf("Error when marshalling request to peer %s: %s\n", conn.GetConnection().RemoteAddr(), err.Error())
        return
    }

    conn.GetConnection().WriteMessage(websocket.TextMessage, data)
}

func (m *DHTManager) createConnection(server string, port uint16, useTls bool) (*websocket.Conn, error) {
    // Connect with timeout
    var timeout uint

    if !c.DecodeMetadata.IsDefined("dht_api", "peer_connection_timeout") {
        timeout = 3 // 5 seconds by default
    } else {
        timeout = c.Conf.DHTApi.PeerConnTimeout
    }

    context, cancel := context.WithTimeout(context.Background(), time.Duration(timeout) * time.Second)
    defer cancel()

    conn, _, err := m.wsDialer.DialContext(context, fmt.Sprintf(
        "ws%s://%s:%d/ws",
        map[bool]string{true: "s", false: ""}[useTls],
        server,
        port,
    ), http.Header{})

    return conn, err
}
