package database

import (
    "gorm.io/gorm"
    "time"
)

// User table model
type User struct {
    Login               string          `gorm:"primaryKey;autoIncrement"`
    Email               string          `gorm:"not null;default:"""`
    IsAdmin             bool            `gorm:"not null"`
    PasswordHash        []byte          `gorm:"not null"`
    PasswordHashSalt    []byte          `gorm:"not null"`

    // Gorm specific field
    CreatedAt           time.Time
    UpdatedAt           time.Time
    DeletedAt           gorm.DeletedAt  `gorm:"index"`
}

// Table name for User structure
func (User) TableName() string {
    return "users"
}

// Games history table model
type GamesHistory struct {
    ID                  int             `gorm:"primaryKey;autoIncrement"`
    GameStartTimestamp  time.Time       `gorm:"not null"`
    GameEndTimestamp    time.Time       `gorm:"not null"`
    IsDraw              bool            `gorm:"not null"`

    // Winner login field
    WinnerLoginKey      string
    WinnerLogin         *User           `gorm:"foreignKey:WinnerLoginKey"`

    // Gorm specific field
    CreatedAt           time.Time
    UpdatedAt           time.Time
    DeletedAt           gorm.DeletedAt  `gorm:"index"`
}

// Table name for GamesHistory structure
func (GamesHistory) TableName() string {
    return "games_history"
}

// Table for storing multiple logins for each game
type PlayerList struct {
    gorm.Model

    LoginKey            string
    Login               *User           `gorm:"foreignKey:LoginKey"`

    GamesHistoryKey     int
    GamesHistory        *GamesHistory   `gorm:"not null;foreignKey:GamesHistoryKey"`
}

func (PlayerList) TableName() string {
    return "player_list"
}

// Table for storing hosts of server instances, detected via handshaking
type DHTHosts struct {
    ServerIdentifier    []byte          `gorm:"primaryKey"`
    SrvLocalIdentifier  []byte          `gorm:"not null"`
    IPAddress           string          `gorm:"not null"`
    Port                uint16          `gorm:"not null"`
    UseTLS              bool            `gorm:"not null"`
    IsPeerStatic        bool            `gorm:"not null"`
    IsPeerConnsStatic   bool            `gorm:"not null"`
    LastHandshake       time.Time

    // Gorm specific field
    CreatedAt           time.Time
    UpdatedAt           time.Time
    DeletedAt           gorm.DeletedAt  `gorm:"index"`
}

func (DHTHosts) TableName() string {
    return "dht_hosts"
}
