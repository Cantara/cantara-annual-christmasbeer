package store

import (
	"context"
	"fmt"
	"github.com/cantara/gober/persistenteventmap"
	"github.com/cantara/gober/stream"
	"go/types"
)

type storeService[pt any] struct {
	beers persistenteventmap.EventMap[Beer, types.Nil]
}

var cryptKey = "MdgKIHmlbRszXjLbS7pXnSBdvl+SR1bSejtpFTQXxro="

func Init[pt any](st stream.Stream[Beer, types.Nil], ctx context.Context) (s storeService[pt], err error) {
	beerMap, err := persistenteventmap.Init[Beer, types.Nil](st, "beer", "0.1.0", func(key string) string {
		return cryptKey
	}, func(b Beer) string {
		return b.ToId()
	}, ctx)
	if err != nil {
		return
	}
	s = storeService[pt]{
		beers: beerMap,
	}
	return
}

func (s storeService[pt]) Set(b Beer) (err error) {
	err = s.beers.Set(b, types.Nil{})
	if err != nil {
		return
	}
	return
}

func (s storeService[pt]) Get(id string) (b Beer, err error) {
	b, _, err = s.beers.Get(id)
	return
}

type Beer struct {
	Name     string  `json:"name"`
	Brand    string  `json:"brand"`
	BrewYear int     `json:"brew_year"`
	ABV      float32 `json:"abv"`
}

func (b Beer) ToId() string {
	return fmt.Sprintf("%s_%s_%d", b.Brand, b.Name, b.BrewYear)
}
