package store

import (
	"context"
	accountTypes "github.com/cantara/cantara-annual-christmasbeer/account/types"
	"github.com/cantara/gober/persistenteventmap"
	"github.com/cantara/gober/store/inmemory"
	"github.com/gofrs/uuid"
	"go/types"
)

type storeService struct {
	users  persistenteventmap.EventMap[accountTypes.Account, types.Nil]
	logins persistenteventmap.EventMap[accountTypes.Login, types.Nil]
}

var cryptKey = "2Q82oDggY6CwBs6QHFu3brYjt8JqFILnn68FDN/eTcU="

func Init(ctx context.Context) (s storeService, err error) {
	store, err := inmemory.Init()
	if err != nil {
		return
	}
	res, err := persistenteventmap.Init[accountTypes.Account, types.Nil](store, "user", "0.1.0", "register_user",
		"create_register_user", "update_register_user", "delete_register_user", func(key string) string {
			return cryptKey
		}, ctx)
	if err != nil {
		return
	}
	los, err := persistenteventmap.Init[accountTypes.Login, types.Nil](store, "login", "0.1.0", "user_login",
		"create_user_login", "update_user_login", "delete_user_login", func(key string) string {
			return cryptKey
		}, ctx)
	if err != nil {
		return
	}
	s = storeService{
		users:  res,
		logins: los,
	}
	return
}

func (s storeService) Register(account accountTypes.Account) (err error) {
	err = s.users.Set(account.Id.String(), account, types.Nil{})
	if err != nil {
		return
	}
	return
}

func (s storeService) Link(username string, login accountTypes.Login) (err error) {
	err = s.logins.Set(username, login, types.Nil{})
	if err != nil {
		return
	}
	return
}

func (s storeService) GetById(id uuid.UUID) (user accountTypes.Account, err error) {
	return s.getUser(id)
}

func (s storeService) GetLogin(username string) (login accountTypes.Login, err error) {
	login, err = s.getLogin(username)
	return
}

func (s storeService) GetByUsername(username string) (user accountTypes.Account, err error) {
	login, err := s.getLogin(username)
	if err != nil {
		return
	}
	return s.getUser(login.AccountId)
}

func (s storeService) getUser(id uuid.UUID) (user accountTypes.Account, err error) {
	user, _, err = s.users.Get(id.String())
	return
}

func (s storeService) getLogin(username string) (login accountTypes.Login, err error) {
	login, _, err = s.logins.Get(username)
	return
}
