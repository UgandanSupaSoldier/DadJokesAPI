package endpoints

import "time"

type JokeViewmodel struct {
	Text       string `validate:"required"`
	Author     *string
	Category   *string
	Rating     *float32 `validate:"max=10"`
	InsertedAt time.Time
}

type TokenViewmodel struct {
	Token     string
	ExpiresAt *time.Time
}
