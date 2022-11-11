package store

import (
	"context"
	"github.com/cantara/gober/persistenteventmap"
	"github.com/cantara/gober/stream"
	"github.com/gofrs/uuid"
	"go/types"
)

type Privilege[t any] struct {
	AccountID uuid.UUID `json:"account_id"`
	Rights    t         `json:"rights"`
}

type storeService[pt any] struct {
	accounts persistenteventmap.EventMap[Privilege[pt], types.Nil]
}

var cryptKey = "2Q82oDggY6CwBs6QHFu3brYjt8JqFILnn68FDN/eTcU="

func Init[pt any](st stream.Stream[Privilege[pt], types.Nil], ctx context.Context) (s storeService[pt], err error) {
	acc, err := persistenteventmap.Init[Privilege[pt], types.Nil](st, "privilege", "0.1.0", func(key string) string {
		return cryptKey
	}, func(p Privilege[pt]) string {
		return p.AccountID.String()
	}, ctx)
	if err != nil {
		return
	}
	s = storeService[pt]{
		accounts: acc,
	}
	return
}

func (s storeService[pt]) Register(accountId uuid.UUID, rights pt) (err error) {
	err = s.accounts.Set(accountId.String(), Privilege[pt]{
		AccountID: accountId,
		Rights:    rights,
	}, types.Nil{})
	if err != nil {
		return
	}
	return
}

func (s storeService[pt]) IsAdmin(id uuid.UUID) bool {
	_, _, err := s.accounts.Get(id.String())
	if err != nil {
		return false
	}
	return true
}

type Beer struct {
	Name     string  `json:"name"`
	Brand    string  `json:"brand"`
	BrewYear int     `json:"brew_year"`
	ABV      float32 `json:"abv"`
}
