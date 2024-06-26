package account

import "github.com/gofrs/uuid"

type AccountRegister struct {
	Username  string `json:"username"`
	Email     string `json:"email"`
	FirstName string `json:"firstname"`
	LastName  string `json:"lastname"`
	Number    string `json:"number"`
	Password  string `json:"password"`
}

type InternalLoggin struct {
	AccountId uuid.UUID `json:"account_id"`
	Username  string    `json:"username"`
	Salt      []byte    `json:"salt"`
	Password  []byte    `json:"password"`
}

type Crypt struct {
	Id  string `json:"id"`
	Key string `json:"key"`
}

type Rights struct {
	Weight float32 `json:"weight"`
	Admin  bool    `json:"admin"`
}
