package tests

import (
	"DadJokesAPI/server"
	"DadJokesAPI/server/database"
	"encoding/json"
	"io"
	"testing"
	"time"

	"github.com/steinfletcher/apitest"
	"github.com/stretchr/testify/assert"
)

func Test_NewToken(t *testing.T) {
	handler := apitest.New().Handler(server.SetupServer())
	request := handler.Get("/v1/token/new")
	result := request.Expect(t).Status(200).End()

	responseBytes, err := io.ReadAll(result.Response.Body)
	assert.Nil(t, err)

	var token database.Token
	err = json.Unmarshal(responseBytes, &token)
	assert.Nil(t, err)

	request = handler.Get("/v1/token/new?expiry=2030-01-01T00:00:00")
	result = request.Expect(t).Status(200).End()

	responseBytes, err = io.ReadAll(result.Response.Body)
	assert.Nil(t, err)

	err = json.Unmarshal(responseBytes, &token)
	assert.Nil(t, err)

	expectedTime, _ := time.Parse("2006-01-02T15:04:05", "2030-01-01T00:00:00")
	assert.Equal(t, token.ExpiresAt, &expectedTime)
}

func Test_Auth(t *testing.T) {

}
