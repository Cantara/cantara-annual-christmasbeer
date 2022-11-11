package store

import (
	"context"
	accountTypes "github.com/cantara/cantara-annual-christmasbeer/account/types"
	"github.com/cantara/gober/persistenteventmap"
	"github.com/cantara/gober/stream"
	"github.com/gofrs/uuid"
	"go/types"
)

type storeService struct {
	users  persistenteventmap.EventMap[accountTypes.Account, types.Nil]
	logins persistenteventmap.EventMap[accountTypes.Login, types.Nil]
}

var cryptKey = "2Q82oDggY6CwBs6QHFu3brYjt8JqFILnn68FDN/eTcU="

func Init(store stream.Persistence, ctx context.Context) (s storeService, err error) {
	acc, err := stream.Init[accountTypes.Account, types.Nil](store, "account", ctx)
	if err != nil {
		return
	}
	lin, err := stream.Init[accountTypes.Login, types.Nil](store, "login", ctx)
	if err != nil {
		return
	}
	res, err := persistenteventmap.Init[accountTypes.Account, types.Nil](acc, "account", "0.1.0", func(key string) string {
		return cryptKey
	}, func(a accountTypes.Account) string {
		return a.Id.String()
	}, ctx)
	if err != nil {
		return
	}
	los, err := persistenteventmap.Init[accountTypes.Login, types.Nil](lin, "login", "0.1.0", func(key string) string {
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
	err = s.users.Set(account, types.Nil{})
	if err != nil {
		return
	}
	return
}

func (s storeService) Link(login accountTypes.Login) (err error) {
	err = s.logins.Set(login, types.Nil{})
	if err != nil {
		return
	}
	return
}

func (s storeService) Accounts() (accounts []accountTypes.Account, err error) {
	s.users.Range(func(_ string, account accountTypes.Account, _ types.Nil) error {
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
	user, _, err = s.users.Get(id.String())
	return
}

func (s storeService) getLogin(username string) (login accountTypes.Login, err error) {
	login, _, err = s.logins.Get(username)
	return
}
