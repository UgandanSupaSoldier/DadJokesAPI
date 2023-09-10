package middleware

import (
	"encoding/json"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	log "github.com/sirupsen/logrus"
)

type requestLog struct {
	Status       int    `json:"status"`
	RequestBody  []byte `json:"request_body"`
	ResponseBody []byte `json:"response_body"`
}

var LogRequest = middleware.BodyDump(
	func(c echo.Context, reqBody []byte, resBody []byte) {
		go func() {
			requestLog := requestLog{
				Status:       c.Response().Status,
				RequestBody:  reqBody,
				ResponseBody: resBody,
			}
			logData, _ := json.Marshal(requestLog)
			log.Info(string(logData))
		}()
	},
)
