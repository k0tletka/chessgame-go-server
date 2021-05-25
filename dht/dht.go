package dht

import (
    "os"
    "log"
    "sync"
    "net/http"
    "path/filepath"
    "io/ioutil"
    "crypto/rand"

    "GoChessgameServer/logger"
)

var (
    // Logger for dht module
    dhtLogger *log.Logger

    // Server identifier
    dhtServerIdentifier [16]byte

    // DHT manager
    DHTMgr *DHTManager
)

func init() {
    // Create logger
    dhtLogger = logger.AddNewLogger("DHTNetwork", os.Stdout, log.LstdFlags | log.Lmsgprefix)
}

// This function called directly from main.go and performs
// server raising and hanshaking with static peers
func InitializeDHT(srvWaitor *sync.WaitGroup, serverChan chan<- *http.Server) {
    // Read or generate server identifier
    generateServerIdentifier()

    // Create new instance of DHT manager
    DHTMgr = CreateNewDHTManager()

    // Start HTTP server
    go InitializeDHTAPIServer(srvWaitor, serverChan)
}


// Function creates file in cache dir and generates identifier for server
func generateServerIdentifier() {
    var cacheDir string
    var err error

    cacheDir, err = os.UserCacheDir()
    if err != nil {
        cacheDir = os.TempDir()
    }

    fpath := filepath.Join(cacheDir, "chessgame_dht")
    _, err = os.Stat(fpath)

    if os.IsNotExist(err) {
        err = os.Mkdir(fpath, 1755)
        if err != nil {
            dhtLogger.Fatalln("Error when creating directory for identificator storing: ", err.Error())
        }
    } else if err != nil {
        dhtLogger.Fatalln("Error when creating directory for identificator storing: ", err.Error())
    }

    file, err := os.OpenFile(filepath.Join(fpath, ".dht_identifier"), os.O_RDWR | os.O_CREATE, 0755)
    if err != nil {
        dhtLogger.Fatalln("Can't open or create file with identifier: ", err.Error())
    }

    // Read identifier from file
    id, err := ioutil.ReadAll(file)
    if err != nil || len(id) == 0 {
        // Generate new server identifier and write it to file
        _, err := rand.Read(dhtServerIdentifier[:])

        if err != nil {
            dhtLogger.Fatalln("Can't generate server identifier: ", err.Error())
        }

        // Ignore errors
        _, _ = file.WriteString(string(dhtServerIdentifier[:]))
        file.Sync()
    }

    copy(dhtServerIdentifier[:], id)
}
