package middleware

import (
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

var RecoverFromPanic = middleware.RecoverWithConfig(middleware.RecoverConfig{
	LogErrorFunc: func(c echo.Context, err error, stack []byte) error {
		log.Errorf("panic error caught: %v", err)
		return nil
	},
})
