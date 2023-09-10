package endpoints

import (
	"DadJokesAPI/server/database"
	"DadJokesAPI/server/responses"
	"math"
	"strconv"

	"github.com/labstack/echo/v4"
)

func RandomJoke(c echo.Context) error {
	joke, err := database.RandomJoke()
	if err != nil {
		return err
	}

	return c.JSON(responses.GenerateResponse(joke, nil))
}

func SearchJokes(c echo.Context) error {
	pageStr := c.QueryParam("page")
	sizeStr := c.QueryParam("page_size")

	page, err := strconv.Atoi(pageStr)
	if err != nil {
		return responses.ErrorWithDetails(responses.InvalidQueryError, "page is required")
	}
	size, err := strconv.Atoi(sizeStr)
	if err != nil {
		return responses.ErrorWithDetails(responses.InvalidQueryError, "page_size is required")
	}

	jokes, err := database.SearchJokes(page, size)
	if err != nil {
		return err
	}
	return c.JSON(responses.GenerateResponse(jokes, nil))
}

func CreateJoke(c echo.Context) error {
	var inputJoke JokeViewmodel
	err := c.Bind(&inputJoke)
	if err != nil {
		return responses.InvalidJSONError
	}

	errors := make(map[string]string)
	if len(inputJoke.Text) == 0 {
		errors["text"] = "text is required"
	}
	if inputJoke.Rating != nil && *inputJoke.Rating > 10 {
		errors["rating"] = "rating must be between 0 and 10"
	}
	if len(errors) > 0 {
		return responses.ErrorWithDetails(responses.InvalidDataError, errors)
	}

	if inputJoke.Rating != nil {
		rating := float32(math.Round(float64(*inputJoke.Rating)*10) / 10)
		inputJoke.Rating = &rating
	}

	joke, err := database.CreateJoke(database.Joke{
		Text:     inputJoke.Text,
		Author:   inputJoke.Author,
		Category: inputJoke.Category,
		Rating:   inputJoke.Rating,
	})
	if err != nil {
		return err
	}

	return c.JSON(responses.GenerateResponse(joke, nil))
}
