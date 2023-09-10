package database

import (
	"DadJokesAPI/server/responses"
	"errors"
	"fmt"
	"time"

	"gorm.io/gorm"
)

func ValidToken(token string) error {
	db, err := connect()
	if err != nil {
		return err
	}

	var rdbToken Token
	err = db.First(&rdbToken, "token = ?", token).Error
	if err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return responses.UnauthorizedError
		} else {
			return fmt.Errorf("failed to get token: %w", err)
		}
	}

	if rdbToken.ExpiresAt != nil && rdbToken.ExpiresAt.Before(time.Now()) {
		return responses.AuthExpiredError
	}
	return nil
}

func CreateToken(token string, expires_at *time.Time) (Token, error) {
	db, err := connect()
	if err != nil {
		return Token{}, err
	}

	rdbToken := Token{
		Token:     token,
		ExpiresAt: expires_at,
	}

	err = db.Create(&rdbToken).Error
	if err != nil {
		return Token{}, fmt.Errorf("failed to create token: %w", err)
	}
	return rdbToken, nil
}
