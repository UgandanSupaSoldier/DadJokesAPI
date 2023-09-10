package tests

import (
	"DadJokesAPI/server"
	"DadJokesAPI/server/database"
	"DadJokesAPI/shared"
	"encoding/json"
	"fmt"
	"io"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

type jokeResponseBody struct {
	Data  database.Joke          `json:"data"`
	Error map[string]interface{} `json:"error"`
}

type jokesResponseBody struct {
	Data  []database.Joke        `json:"data"`
	Error map[string]interface{} `json:"error"`
}

func Test_RandomJoke(t *testing.T) {
	shared.SeConfigtPath("../config.toml")

	handler := apitest.New().Handler(server.SetupServer())
	request := handler.Get("/v1/joke/random")
	result := request.Expect(t).Status(200).End()

	responseBytes, err := io.ReadAll(result.Response.Body)
	assert.Nil(t, err)

	var joke jokeResponseBody
	err = json.Unmarshal(responseBytes, &joke)
	assert.NotEmpty(t, joke.Data.ID)
	assert.Nil(t, err)
}

func Test_PagedJokes(t *testing.T) {
	shared.SeConfigtPath("../config.toml")

	handler := apitest.New().Handler(server.SetupServer())
	request := handler.Get("/v1/joke/search").QueryParams(map[string]string{
		"page":      "20",
		"page_size": "3",
	})

	jokes := []database.Joke{
		{
			ID:       61,
			Text:     "Did you hear about the circus fire?<> It was in tents!",
			Author:   nil,
			Category: nil,
			Rating:   nil,
		},
		{
			ID:       62,
			Text:     "Don't trust atoms.<> They make up everything!",
			Author:   nil,
			Category: nil,
			Rating:   nil,
		},
		{
			ID:       63,
			Text:     "How many tickles does it take to make an octopus laugh? <>Ten-tickles.",
			Author:   nil,
			Category: nil,
			Rating:   nil,
		},
	}
	responseBody, err := json.Marshal(jokesResponseBody{Data: jokes, Error: map[string]interface{}{}})
	assert.Nil(t, err)

	request.Expect(t).Status(200).Body(string(responseBody)).End()
}

func Test_CreateJoke(t *testing.T) {
	shared.SeConfigtPath("../config.toml")

	handler := apitest.New().Handler(server.SetupServer())
	request := handler.Get("/v1/token/new")
	result := request.Expect(t).Status(200).End()

	responseBytes, err := io.ReadAll(result.Response.Body)
	assert.Nil(t, err)

	var tokenResponseBody tokenResponseBody
	err = json.Unmarshal(responseBytes, &tokenResponseBody)
	assert.Nil(t, err)

	jokeInput := map[string]interface{}{
		"text":     "Did you hear about the circus fire?<> It was in tents!",
		"author":   "John Doe",
		"category": "Programming",
		"rating":   7.45,
	}

	request = handler.Post("/v1/joke/create").JSON(jokeInput).Header("Authorization", fmt.Sprintf("bearer %s", tokenResponseBody.Data.Token))
	result = request.Expect(t).Status(200).End()

	responseBytes, err = io.ReadAll(result.Response.Body)
	assert.Nil(t, err)

	var jokeResponseBody jokeResponseBody
	err = json.Unmarshal(responseBytes, &jokeResponseBody)
	assert.Nil(t, err)

	assert.Equal(t, jokeInput["text"], jokeResponseBody.Data.Text)
	assert.Equal(t, jokeInput["author"], *jokeResponseBody.Data.Author)
	assert.Equal(t, jokeInput["category"], *jokeResponseBody.Data.Category)
	assert.Equal(t, float32(7.4), *jokeResponseBody.Data.Rating)
}
