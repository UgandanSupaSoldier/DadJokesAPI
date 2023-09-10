package database

import (
	"DadJokesAPI/shared"
	"fmt"
)

func CreateLog(log Log) error {
	db, err := shared.Connect()
	if err != nil {
		return err
	}

	err = db.Create(&log).Error
	if err != nil {
		return fmt.Errorf("failed to create joke: %w", err)
	}
	return nil
}
