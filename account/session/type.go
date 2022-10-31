package session

import "time"

type AccessTokenSession struct {
	AccessToken AccessToken `json:"access_token"`
	AccountId   string      `json:"account_id"`
}

type AccessToken struct {
	Token     string    `json:"token"`
	ExpiresIn int       `json:"expires_in"`
	Expires   time.Time `json:"expires"`
}
