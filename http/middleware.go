package http

import (
	"mysql/app/apperr"

	"github.com/labstack/echo/v4"
)

// RecoverPanicMiddleware is the middleware for handling panics.
func (s *ServerAPI) RecoverPanicMiddleware(next echo.HandlerFunc) echo.HandlerFunc {
	return func(c echo.Context) error {

		defer func() {
			if r := recover(); r != nil {
				err := apperr.Errorf(apperr.EUNKNOWN, "panic: %s", r)
				ErrorResponseJSON(c, err, nil)
			}
		}()

		return next(c)
	}
}
