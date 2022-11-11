package account

import (
	"github.com/cantara/cantara-annual-christmasbeer/account/privilege"
	"github.com/cantara/cantara-annual-christmasbeer/account/session"
	"github.com/cantara/cantara-annual-christmasbeer/account/types"
	"github.com/gofrs/uuid"
)

type StoreService interface {
	Register(account types.Account) (err error)
	Link(login types.Login) (err error)
	Accounts() (accounts []types.Account, err error)
	GetById(id uuid.UUID) (user types.Account, err error)
	GetLogin(username string) (login types.Login, err error)
	GetByUsername(username string) (user types.Account, err error)
}

type SessionService interface {
	Create(accountId uuid.UUID) (accessToken session.AccessToken, err error)
	Renew(token string) (accessToken session.AccessToken, err error)
	Validate(token string) (accessToken session.AccessToken, accountId uuid.UUID, err error)
}

type AdminService interface {
	Register(accountId uuid.UUID, rights Rights) (err error)
	IsAdmin(accountId uuid.UUID) (p privilege.Privilege[Rights], err error)
}
