package middleware

import (
	"DadJokesAPI/server/database"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

var LogRequest = middleware.BodyDump(
	func(c echo.Context, reqBody []byte, resBody []byte) {
		go func() {
			dbLog := database.Log{
				Time:         time.Now(),
				RequestUrl:   c.Request().URL.String(),
				RequestBody:  string(reqBody),
				ResponseBody: string(resBody),
				ResponseCode: c.Response().Status,
			}

			err := database.CreateLog(dbLog)
			if err != nil {
				log.Errorf("failed to create log %v: %v", dbLog, err)
			}
		}()
	},
)
