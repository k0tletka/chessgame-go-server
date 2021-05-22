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
    ID                  int         `gorm:"primaryKey;autoIncrement"`
    GameStartTimestamp  time.Time   `gorm:"not null"`
    GameEndTimestamp    time.Time   `gorm:"not null"`
    IsDraw              bool        `gorm:"not null"`

    // Winner login field
    WinnerLoginKey      string
    WinnerLogin         *User        `gorm:"foreignKey:WinnerLoginKey"`

    // Player one field
    PlayerOneLoginKey   string
    PlayerOneLogin      *User        `gorm:"not null;foreignKey:PlayerOneLoginKey"`

    // Player two field
    PlayerTwoLoginKey   string
    PlayerTwoLogin      *User        `gorm:"not null;foreignKey:PlayerTwoLoginKey"`

    // Gorm specific field
    CreatedAt           time.Time
    UpdatedAt           time.Time
    DeletedAt           gorm.DeletedAt  `gorm:"index"`
}

// Table name for GamesHistory structure
func (GamesHistory) TableName() string {
    return "games_history"
}
