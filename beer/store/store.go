package store

import (
	"context"
	"fmt"

	"github.com/cantara/bragi/sbragi"
	"github.com/cantara/gober/persistenteventmap"
	"github.com/cantara/gober/stream"
	"github.com/cantara/gober/stream/event"
	"github.com/cantara/gober/stream/event/store"
)

type storeService[pt any] struct {
	beers persistenteventmap.EventMap[Beer]
}

var cryptKey = "MdgKIHmlbRszXjLbS7pXnSBdvl+SR1bSejtpFTQXxro="

func Init[pt any](st stream.Stream, ctx context.Context) (s storeService[pt], err error) {
	beerMap, err := persistenteventmap.Init[Beer](st, "beer", "0.1.0", stream.StaticProvider(sbragi.RedactedString(cryptKey)), func(b Beer) string {
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
	err = s.beers.Set(b)
	if err != nil {
		return
	}
	return
}

func (s storeService[pt]) Get(id string) (b Beer, err error) {
	b, err = s.beers.Get(id)
	return
}
func (s storeService[pt]) Stream(ctx context.Context) (out <-chan event.Event[Beer], err error) {
	return s.beers.Stream(event.AllTypes(), store.STREAM_START, stream.ReadDataType("beer"), ctx)
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
