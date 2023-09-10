package endpoints

import "time"

type JokeViewmodel struct {
	Text     string
	Author   *string
	Category *string
	Rating   *float32
}

type TokenViewmodel struct {
	Token     string
	ExpiresAt *time.Time
}
