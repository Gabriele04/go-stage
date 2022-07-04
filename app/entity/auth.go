package entity

import (
	"fmt"
	"time"

	"github.com/music-gang/music-gang-api/app/apperr"
	"gopkg.in/guregu/null.v4"
)

// Authentication providers.
const (
	AuthSourceLocal  = "local"
	AuthSourceGitHub = "github"
)

// Auth represents a set of OAuth credentials. These are linked to a User so a
// single user could authenticate through multiple providers.
//
// The authentication system links users by email address, however, some GitHub
// users don't provide their email publicly so we may not be able to link them
// by email address.
type Auth struct {
	ID       int64       `json:"id"`
	CityID   int64       `json:"city_id"`
	Source   string      `json:"source"`
	SourceID null.String `json:"-"`

	AccessToken  null.String `json:"-"`
	RefreshToken null.String `json:"-"`
	Expiry       null.Time   `json:"-"`

	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`

	City *City `json:"city"`
}

// Auths represents a list of Auth objects.
type Auths []*Auth

// AuthCityOptions represents the options for a user during fetching auth service.
type AuthCityOptions struct {
	AuthCode   *string
	Source     *string
	CityParams *CityParams
}

// UserParams represents the parameters for a city authentication.
type CityParams struct {
	Name       *string
	Population *int
}

// CanAuthBeDeleted returns if the passed authentication source can be deleted.
func CanAuthBeDeleted(auth *Auth) bool {

	var canBeDeleted bool

	switch auth.Source {

	case AuthSourceGitHub:
		canBeDeleted = true

	default:
		canBeDeleted = false
	}

	return canBeDeleted
}

// IsSourceIDRequired returns if the authentication source requires a source ID.
func IsSourceIDRequired(source string) bool {
	switch source {
	case AuthSourceGitHub:
		return true
	default:
		return false
	}
}

// AvatarURL returns a URL to the avatar image hosted by the authentication source.
// Returns an empty string if the authentication source is invalid.
func (a *Auth) AvatarURL(size int) string {
	switch a.Source {
	case AuthSourceGitHub:
		return fmt.Sprintf("https://avatars1.githubusercontent.com/u/%s?s=%d", a.SourceID.String, size)
	default:
		return ""
	}
}

// Validate validates the Auth object and returns an error if it's invalid.
// This can be used from any method that accepts a Auth object to do basic checks.
func (a *Auth) Validate() error {

	if a.CityID == 0 {

		return apperr.Errorf(apperr.EINVALID, "User is required")

	} else if a.Source == "" {

		return apperr.Errorf(apperr.EINVALID, "Source is required")

	} else if IsSourceIDRequired(a.Source) && a.SourceID.String == "" {

		return apperr.Errorf(apperr.EINVALID, "Source ID is required")

	} else if a.SourceID.Valid && a.SourceID.String == "" {

		return apperr.Errorf(apperr.EINVALID, "Source ID cannot be empty if provided")

	} else if a.AccessToken.Valid && a.AccessToken.String == "" {

		return apperr.Errorf(apperr.EINVALID, "Access token cannot be empty if provided")

	} else if a.RefreshToken.Valid && a.RefreshToken.String == "" {

		return apperr.Errorf(apperr.EINVALID, "Refresh token cannot be empty if provided")
	}

	return nil
}
