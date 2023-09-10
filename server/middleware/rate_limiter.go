package middleware

import (
	"DadJokesAPI/server/responses"
	"DadJokesAPI/shared"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
	"golang.org/x/time/rate"
)

var (
	rateLimiter *echo.MiddlewareFunc
)

func newMemoryStore() middleware.RateLimiterStore {
	rateLimit := shared.GetIntDef("rate_limiter.rate", 1)
	burstLimit := shared.GetIntDef("rate_limiter.burst", 10)
	timeout := shared.GetIntDef("rate_limiter.timeout", 180)

	return middleware.NewRateLimiterMemoryStoreWithConfig(
		middleware.RateLimiterMemoryStoreConfig{
			Rate: rate.Limit(rateLimit), Burst: burstLimit, ExpiresIn: time.Duration(timeout) * time.Second,
		},
	)
}

func getRateLimiter() *echo.MiddlewareFunc {
	if rateLimiter == nil {
		newRateLimiter := middleware.RateLimiterWithConfig(
			middleware.RateLimiterConfig{
				Store: newMemoryStore(),
				DenyHandler: func(c echo.Context, identifier string, err error) error {
					return c.JSON(responses.GenerateResponse(nil, responses.TooManyRequestsError))
				},
			},
		)
		rateLimiter = &newRateLimiter
	}
	return rateLimiter
}

func DynamicRateLimiter() echo.MiddlewareFunc {
	return func(next echo.HandlerFunc) echo.HandlerFunc {
		return func(c echo.Context) error {
			if shared.GetBoolDef("rate_limiter.enabled", false) {
				if rl := getRateLimiter(); rl != nil {
					return (*rl)(next)(c)
				}
			}
				return next(c)
		}
	}
}
