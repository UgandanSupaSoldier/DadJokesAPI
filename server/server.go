package server

import (
	"DadJokesAPI/server/endpoints"
	"DadJokesAPI/server/middleware"
	"DadJokesAPI/shared"
	"fmt"

	"github.com/labstack/echo/v4"
)

func Run() error {
	server := SetupServer()

	port := shared.GetIntDef("server.port", 8080)
	err := server.Start(fmt.Sprintf("0.0.0.0:%d", port))
	if err != nil {
		return fmt.Errorf("failed to start server: %w", err)
	}
	return nil
}

func SetupServer() *echo.Echo {
	server := echo.New()

	setupMiddleware(server)
	setupRoutes(server)

	return server
}

func setupMiddleware(server *echo.Echo) {
	server.Use(middleware.LogRequest)
	server.Use(middleware.RecoverFromPanic)
	server.Use(middleware.DynamicRateLimiter())

	server.HTTPErrorHandler = middleware.HTTPErrorHandler
}

func setupRoutes(server *echo.Echo) {
	v1api := server.Group("/v1")

	token := v1api.Group("/token")
	token.GET("/new", endpoints.NewToken)

	joke := v1api.Group("/joke")
	joke.GET("/random", endpoints.RandomJoke)
	joke.GET("/search", endpoints.SearchJokes)
	joke.POST("/create", endpoints.CreateJoke, middleware.AuthUser)
}
