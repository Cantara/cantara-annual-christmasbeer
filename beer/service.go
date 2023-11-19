package beer

import (
	"context"

	"github.com/cantara/cantara-annual-christmasbeer/beer/store"
	"github.com/cantara/gober/stream"
	"github.com/cantara/gober/stream/event"
)

type Store interface {
	Set(b store.Beer) (err error)
	Get(id string) (b store.Beer, err error)
	Stream(ctx context.Context) (out <-chan event.Event[store.Beer], err error)
	Range(f func(key string, data store.Beer) error)
}

type service struct {
	store Store
}

func InitService(st stream.Stream, ctx context.Context) (s service, err error) {
	beerStore, err := store.Init[store.Beer](st, ctx)
	if err != nil {
		return
	}
	s = service{
		store: beerStore,
	}
	return
}

func (s service) Get(id string) (b store.Beer, err error) {
	b, err = s.store.Get(id)
	return
}

func (s service) Range(f func(key string, data store.Beer) error) {
	s.store.Range(f)
}

func (s service) Register(b store.Beer) (err error) {
	return s.store.Set(b)
}

func (s service) BeerStream(ctx context.Context) (out <-chan event.Event[store.Beer], err error) {
	out, err = s.store.Stream(ctx)
	return
}
