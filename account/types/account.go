package types

import (
	"github.com/gofrs/uuid"
	"time"
)

type Codes struct {
	Email  string `json:"email"`
	Mobile string `json:"mobile"`
	MFA    string `json:"mfa"`
}

type Account struct {
	Id                uuid.UUID `json:"id"`
	Email             string    `json:"email"`
	FirstName         string    `json:"firstname"`
	LastName          string    `json:"lastname"`
	Number            string    `json:"number"`
	VerificationCodes Codes     `json:"verification_codes"`
	UpdatedAt         time.Time `json:"updated_at"`
	CreatedAt         time.Time `json:"created_at"`
}
