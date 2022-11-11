package types

import (
	"github.com/gofrs/uuid"
)

type LoginType int

var loginTypeStrings = []string{"INVALID", "INTERNAL", "OAUTH", "OIDC"}

const (
	INVALID LoginType = iota
	INTERNAL
	OAUTH
	OIDC
)

func (l LoginType) String() string {
	return loginTypeStrings[l]
}

func LoginTypeFromString(loginType string) LoginType {
	switch loginType {
	case loginTypeStrings[0]:
		return INVALID
	case loginTypeStrings[1]:
		return INTERNAL
	case loginTypeStrings[2]:
		return OAUTH
	case loginTypeStrings[3]:
		return OIDC
	default:
		return INVALID
	}
}

type Login struct {
	Id        string    `json:"id"`
	AccountId uuid.UUID `json:"account_id"`
	Type      LoginType `json:"login_type"`
	Data      []byte    `json:"data"`
}
