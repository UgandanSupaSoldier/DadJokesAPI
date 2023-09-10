package middleware

import (
	"DadJokesAPI/server/responses"
	"errors"

	"github.com/labstack/echo/v4"

	log "github.com/sirupsen/logrus"
)

func HTTPErrorHandler(err error, c echo.Context) {
	if !c.Response().Committed {
		var handledError *responses.ErrorType
		internalServerError := &responses.ErrorType{
			Code: 500, ResponseError: responses.ResponseError{Type: "internal_server_error",
				Message: "The server encountered an internal error and was unable to complete your request", Details: nil},
		}

		if errors.Is(err, echo.ErrNotFound) || errors.Is(err, echo.ErrMethodNotAllowed) {
			handledError = responses.NotFoundError
		} else {
			if ok := errors.As(err, &handledError); !ok {
				log.Errorf("unhandled server error: %v", err)
				handledError = internalServerError
			} else if handledError == nil {
				log.Errorf("handled server error with nil value: %v", err)
				handledError = internalServerError
			}
		}

		err = c.JSON(responses.GenerateResponse(nil, handledError))
		if err != nil {
			log.Errorf("failed to generate JSON error response: %v", err)
		}
	}
}
