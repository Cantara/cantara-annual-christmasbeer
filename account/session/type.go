package session

import (
	"github.com/gofrs/uuid"
	"time"
)

type AccessTokenSession struct {
	AccessToken AccessToken `json:"access_token"`
	AccountId   uuid.UUID   `json:"account_id"`
}

type AccessToken struct {
	Token     string    `json:"token"`
	ExpiresIn int       `json:"expires_in"`
	Expires   time.Time `json:"expires"`
}
