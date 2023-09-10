package tests

import (
	"DadJokesAPI/server"
	"DadJokesAPI/server/database"
	"DadJokesAPI/shared"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

type tokenResponseBody struct {
	Data  database.Token         `json:"data"`
	Error map[string]interface{} `json:"error"`
}

type errorResponse struct {
	Data  map[string]interface{} `json:"data"`
	Error map[string]interface{} `json:"error"`
}

func Test_NewToken(t *testing.T) {
	shared.SeConfigtPath("../config.toml")

	handler := apitest.New().Handler(server.SetupServer())
	request := handler.Get("/v1/token/new")
	result := request.Expect(t).Status(200).End()

	responseBytes, err := io.ReadAll(result.Response.Body)
	assert.Nil(t, err)

	var token tokenResponseBody
	err = json.Unmarshal(responseBytes, &token)
	assert.NotEmpty(t, token.Data.Token)
	assert.Nil(t, err)

	request = handler.Get("/v1/token/new").QueryParams(map[string]string{"expiry": "2030-01-01T16:00:30"})
	result = request.Expect(t).Status(200).End()

	responseBytes, err = io.ReadAll(result.Response.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(responseBytes, &token)
	assert.Nil(t, err)

	expectedTime, _ := time.Parse("2006-01-02T15:04:05", "2030-01-01T16:00:30")
	assert.Equal(t, token.Data.ExpiresAt, &expectedTime)
}

func Test_Auth(t *testing.T) {
	shared.SeConfigtPath("../config.toml")

	msg1 := "The request is missing an authorization header"
	msg2 := "The authorization header is malformed or invalid"
	msg3 := "You are not authorized to access this resource"

	resp1, err := json.Marshal(errorResponse{
		Data:  map[string]interface{}{},
		Error: map[string]interface{}{"type": "unauthorized", "message": msg1},
	})
	assert.Nil(t, err)
	resp2, err := json.Marshal(errorResponse{
		Data:  map[string]interface{}{},
		Error: map[string]interface{}{"type": "unauthorized", "message": msg2},
	})
	assert.Nil(t, err)
	resp3, err := json.Marshal(errorResponse{
		Data:  map[string]interface{}{},
		Error: map[string]interface{}{"type": "unauthorized", "message": msg3},
	})
	assert.Nil(t, err)

	handler := apitest.New().Handler(server.SetupServer())
	request := handler.Post("/v1/joke/create")
	request.Expect(t).Status(401).Body(string(resp1)).End()
	request = handler.Post("/v1/joke/create").Header("Authorization", "1234")
	request.Expect(t).Status(401).Body(string(resp2)).End()
	handler = apitest.New().Handler(server.SetupServer())
	request = handler.Post("/v1/joke/create").Header("Authorization", "bearer 1234")
	request.Expect(t).Status(401).Body(string(resp3)).End()
}
