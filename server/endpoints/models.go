package endpoints

import "time"

type JokeViewmodel struct {
	Text       string `validate:"required"`
	Author     *string
	Category   *string
	Rating     *float32 `validate:"max=5"`
	InsertedAt time.Time
	Tags       []string
}

type TokenViewmodel struct {
	Token     string
	ExpiresAt *time.Time
}
