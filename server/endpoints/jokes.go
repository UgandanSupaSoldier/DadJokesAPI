package endpoints

import (
	"DadJokesAPI/server/database"
	"DadJokesAPI/server/responses"
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

	err = database.CreateJoke(database.Joke{
		Text:     inputJoke.Text,
		Author:   inputJoke.Author,
		Category: inputJoke.Category,
		Rating:   inputJoke.Rating,
		Tags:     inputJoke.Tags,
	})
	if err != nil {
		return err
	}

	return c.JSON(responses.GenerateResponse(nil, nil))
}