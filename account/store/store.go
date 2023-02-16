package store

import (
	"context"
	accountTypes "github.com/cantara/cantara-annual-christmasbeer/account/types"
	"github.com/cantara/gober/persistenteventmap"
	"github.com/cantara/gober/stream"
	"github.com/gofrs/uuid"
)

type storeService struct {
	users  persistenteventmap.EventMap[accountTypes.Account]
	logins persistenteventmap.EventMap[accountTypes.Login]
}

var cryptKey = "2Q82oDggY6CwBs6QHFu3brYjt8JqFILnn68FDN/eTcU="

func Init(store stream.Stream, ctx context.Context) (s storeService, err error) {
	res, err := persistenteventmap.Init[accountTypes.Account](store, "account", "0.1.0", func(key string) string {
		return cryptKey
	}, func(a accountTypes.Account) string {
		return a.Id.String()
	}, ctx)
	if err != nil {
		return
	}
	los, err := persistenteventmap.Init[accountTypes.Login](store, "login", "0.1.0", func(key string) string {
		return cryptKey
	}, func(l accountTypes.Login) string {
		return l.Id
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
	err = s.users.Set(account)
	if err != nil {
		return
	}
	return
}

func (s storeService) Link(login accountTypes.Login) (err error) {
	err = s.logins.Set(login)
	if err != nil {
		return
	}
	return
}

func (s storeService) Accounts() (accounts []accountTypes.Account, err error) {
	s.users.Range(func(_ string, account accountTypes.Account) error {
		accounts = append(accounts, account)
		return nil
	})
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
	user, err = s.users.Get(id.String())
	return
}

func (s storeService) getLogin(username string) (login accountTypes.Login, err error) {
	login, err = s.logins.Get(username)
	return
}
