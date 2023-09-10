package database

import "time"

type Token struct {
	ID        int `gorm:"primaryKey,column:id"`
	Token     string
	ExpiresAt *time.Time
}

type Joke struct {
	ID         int `gorm:"primaryKey,column:id"`
	Text       string
	Author     *string
	Category   *string
	Rating     *float32
	InsertedAt time.Time
	Tags       []string
}
