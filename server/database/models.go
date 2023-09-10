package database

import "time"

type Token struct {
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
}

type Log struct {
	Time         time.Time
	RequestUrl   string
	RequestBody  string
	ResponseBody string
	ResponseCode int
}
