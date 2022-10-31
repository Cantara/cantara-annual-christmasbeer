package oidc

import (
	"fmt"
	"time"
)

type Login struct {
	AccountId             string    `json:"account_id"`
	Username              string    `json:"username"`
	ExternalUsername      string    `json:"external_username"`
	ExternalEmail         string    `json:"external_email"`
	ExternalEmailVerified bool      `json:"external_email_verified"`
	Provider              Provider  `json:"provider"`
	AccessToken           string    `json:"access_token"`
	RefreshToken          string    `json:"refresh_token"`
	ExpiresAt             time.Time `json:"expires_at"`
}

type Provider string

const (
	TWITCH  Provider = "twitch"
	YOUTUBE Provider = "youtube"
)

var INVALID_PROVIDER_ERROR = fmt.Errorf("INVALID PROVIDER")
