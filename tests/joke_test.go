package tests

import (
	"DadJokesAPI/server"
	"DadJokesAPI/server/database"
	"encoding/json"
	"io"
	"testing"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func Test_RandomJoke(t *testing.T) {
	handler := apitest.New().Handler(server.SetupServer())
	request := handler.Get("/v1/joke/random")
	result := request.Expect(t).Status(200).End()

	responseBytes, err := io.ReadAll(result.Response.Body)
	assert.Nil(t, err)

	var joke database.Joke
	err = json.Unmarshal(responseBytes, &joke)
	assert.Nil(t, err)
}

func Test_PagedJokes(t *testing.T) {
	handler := apitest.New().Handler(server.SetupServer())
	request := handler.Get("/v1/joke/search?page=20&page_size=3")

	requestBody, err := json.Marshal([]database.Joke{
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
	})
	assert.Nil(t, err)

	request.Expect(t).Status(200).Body(string(requestBody)).End()
}

func Test_CreateJoke(t *testing.T) {
	handler := apitest.New().Handler(server.SetupServer())
	request := handler.Get("/v1/joke/create").JSON()
	result := request.Expect(t).Status(200).End()

	responseBytes, err := io.ReadAll(result.Response.Body)
	assert.Nil(t, err)

	var joke database.Joke
	err = json.Unmarshal(responseBytes, &joke)
	assert.Nil(t, err)
}
