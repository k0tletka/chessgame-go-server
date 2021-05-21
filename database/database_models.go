package database

import (
    "gorm.io/gorm"
    "time"
)

// User table model
type User struct {
    Login               string          `gorm:"primaryKey;autoIncrement"`
    Email               string          `gorm:"not null"`
    IsAdmin             bool            `gorm:"not null"`
    PasswordHash        []byte          `gorm:"not null"`
    PasswordHashSalt    []byte          `gorm:"not null"`

    // Gorm specific field
    CreatedAt           time.Time
    UpdatedAt           time.Time
    DeletedAt           gorm.DeletedAt  `gorm:"index"`
}

// Games history table model
type GamesHistory struct {
    ID                  int         `gorm:"primaryKey;autoIncrement"`
    GameStartTimestamp  time.Time   `gorm:"not null"`
    GameEndTimestamp    time.Time   `gorm:"not null"`
    IsDraw              bool        `gorm:"not null"`
    WinnerLogin         User
    PlayerOneLogin      User        `gorm:"not null"`
    PlayerTwoLogin      User        `gorm:"not null"`

    // Gorm specific field
    CreatedAt           time.Time
    UpdatedAt           time.Time
    DeletedAt           gorm.DeletedAt  `gorm:"index"`
}
