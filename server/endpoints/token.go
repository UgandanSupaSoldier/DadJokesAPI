package endpoints

import (
	"DadJokesAPI/server/database"
	"DadJokesAPI/server/responses"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func NewToken(c echo.Context) error {
	var expiry *time.Time

	expiryString := c.QueryParam("expiry")
	if len(expiryString) > 0 {
		expiryTime, err := time.Parse(time.RFC3339, expiryString)
		if err != nil {
			return responses.ErrorWithDetails(responses.InvalidQueryError, "expiry is not in valid RFC3339 format")
		}
		if expiryTime.Before(time.Now()) {
			return responses.ErrorWithDetails(responses.InvalidQueryError, "expiry may only contain future dates")
		}
		expiry = &expiryTime
	}

	newToken, err := database.CreateToken(uuid.New().String(), expiry)
	if err != nil {
		return err
	}

	return c.JSON(responses.GenerateResponse(TokenViewmodel{
		Token:     newToken.Token,
		ExpiresAt: newToken.ExpiresAt,
	}, nil))
}
