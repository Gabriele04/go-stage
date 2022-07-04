package http

import (
	"mysql/app/apperr"
	"mysql/app/entity"

	"github.com/labstack/echo/v4"
)

const (
	claimsContextParam = "claims"
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

func AuthCity(c echo.Context) (*entity.City, error) {

	if claims, ok := c.Get(claimsContextParam).(*entity.AppClaims); ok {
		return claims.Auth.City, nil
	}

	return nil, apperr.Errorf(apperr.EUNAUTHORIZED, "no auth city found in context")
}
