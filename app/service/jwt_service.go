package service

import (
	"context"

	"mysql/app/entity"
)

// JWTService is an interface for JWT service.
// It is used to generate and validate JWT tokens.
type JWTService interface {
	// Exchange a auth entity for a JWT token pair.
	Exchange(ctx context.Context, auth *entity.User) (*entity.TokenPair, error)

	// Parse a JWT token and return the associated claims.
	Parse(ctx context.Context, token string) (*entity.AppClaims, error)
}
