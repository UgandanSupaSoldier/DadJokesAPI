package middleware

import (
	"DadJokesAPI/server/database"
	"DadJokesAPI/server/responses"
	"strings"

	"github.com/labstack/echo/v4"
)

func AuthUser(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {
		authHeader := c.Request().Header.Get("Authorization")
		if len(authHeader) == 0 {
			return responses.AuthMissingError
		}
		authHeaderParts := strings.Split(authHeader, " ")
		if len(authHeaderParts) != 2 || strings.ToLower(authHeaderParts[0]) != "bearer" {
			return responses.AuthInvalidError
		}

		err := database.ValidToken(authHeaderParts[1])
		if err != nil {
			return err
		} else {
			return next(c)
		}
	}
}
