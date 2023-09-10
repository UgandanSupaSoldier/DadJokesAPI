package database

import (
	"DadJokesAPI/shared"
	"fmt"
	"time"
)

func RandomJoke() (Joke, error) {
	db, err := shared.Connect()
	if err != nil {
		return Joke{}, err
	}

	var rdbJoke Joke
	err = db.Order("RANDOM()").First(&rdbJoke).Error
	if err != nil {
		return Joke{}, fmt.Errorf("failed to get random joke: %w", err)
	}
	return rdbJoke, nil
}

func SearchJokes(page int, size int) ([]Joke, error) {
	db, err := shared.Connect()
	if err != nil {
		return nil, err
	}

	var rdbJokes []Joke
	err = db.Order("id").Limit(size).Offset(page * size).Find(&rdbJokes).Error
	if err != nil {
		return nil, fmt.Errorf("failed to get paged jokes: %w", err)
	}
	return rdbJokes, nil
}

func CreateJoke(joke Joke) error {
	db, err := shared.Connect()
	if err != nil {
		return err
	}

	joke.InsertedAt = time.Now()
	err = db.Create(&joke).Error
	if err != nil {
		return fmt.Errorf("failed to create joke: %w", err)
	}
	return nil
}
