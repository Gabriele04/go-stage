package jwt

import (
	"context"
	"mysql/app/apperr"
	"mysql/app/entity"
	"mysql/app/service"

	"github.com/golang-jwt/jwt"
)

const (
	accessTokenExpiration = entity.AccessTokenExpiration
)

var _ service.JWTService = (*JWTService)(nil)

type JWTService struct {
	Secret string
}

// Exchange implements service.JWTService
func (s *JWTService) Exchange(ctx context.Context, user *entity.User) (*entity.TokenPair, error) {
	select {
	case <-ctx.Done():
		return nil, apperr.Errorf(apperr.EINTERNAL, "context cancelled")
	default:
		accessTokenclaims := entity.NewAppClaims(user, accessTokenExpiration)

		accessToken := jwt.NewWithClaims(jwt.SigningMethodHS256, accessTokenclaims)

		accessTokenString, err := accessToken.SignedString([]byte(s.Secret))
		if err != nil {
			return nil, apperr.Errorf(apperr.EINTERNAL, "failed to sign access token: %v", err)
		}

		return &entity.TokenPair{
			AccessToken:  accessTokenString,
			TokenType:    "Bearer",
			Expiry:       accessTokenclaims.ExpiresAt,
		}, nil
	}
}

// Parse implements service.JWTService
func (s *JWTService) Parse(ctx context.Context, token string) (*entity.AppClaims, error) {
	select {
	case <-ctx.Done():
		return nil, apperr.Errorf(apperr.EINTERNAL, "context cancelled")
	default:
		t, err := jwt.ParseWithClaims(token, &entity.AppClaims{}, func(token *jwt.Token) (interface{}, error) {
			if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
				return nil, apperr.Errorf(apperr.EINTERNAL, "unexpected signing method: %v", token.Header["alg"])
			}

			return []byte(s.Secret), nil
		})
		if err != nil {
			return nil, apperr.Errorf(apperr.EINTERNAL, "failed to parse token: %v", err)
		}

		claims, ok := t.Claims.(*entity.AppClaims)
		if !ok {
			return nil, apperr.Errorf(apperr.EINTERNAL, "failed to parse claims")
		}

		if !t.Valid {
			return nil, apperr.Errorf(apperr.EUNAUTHORIZED, "invalid token")
		}

		return claims, nil
	}
}
