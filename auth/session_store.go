package auth

import (
    "errors"
    "sync"
)

var (
    // Errors
    SessionRegisteredAlready = errors.New("auth: Session for this login has been registered already")
    NoSuchSession = errors.New("auth: There is no session for given login")

    // Session store. Stores login, jwt and IP addresses for every user
    SessionStore = &SessionStoreType{
        sessions: make(map[string]*SessionInformation),
        mutex: &sync.RWMutex{},
    }
)

// Session types. Stores login name and other authetication and session-related data
type SessionStoreType struct {
    sessions map[string]*SessionInformation
    mutex *sync.RWMutex
}

// Create a new session for given login
func (s *SessionStoreType) CreateNewSession(login string, sinfo *SessionInformation) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    if _, nok := s.sessions[login]; nok && s.sessions[login] != nil {
        return SessionRegisteredAlready
    }

    sinfo.GenerateKey()
    s.sessions[login] = sinfo

    return nil
}

// Terminate session (delete session from map)
func (s *SessionStoreType) TerminateSession(login string) error {
    s.mutex.Lock()
    defer s.mutex.Unlock()

    if _, ok := s.sessions[login]; !ok {
        return NoSuchSession
    }

    if conn := s.sessions[login].WSConnection; conn != nil {
        conn.CloseConnection("Session closed")
    }

    s.sessions[login] = nil
    return nil
}

// Get session
func (s *SessionStoreType) GetSession(login string) (*SessionInformation, error) {
    if s.sessions[login] == nil {
        return nil, NoSuchSession
    }

    return s.sessions[login], nil
}
