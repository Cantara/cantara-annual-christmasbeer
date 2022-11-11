package beer

import (
	"context"
	log "github.com/cantara/bragi"
	"github.com/cantara/cantara-annual-christmasbeer/beer/store"
	streamStore "github.com/cantara/gober/store"
	"github.com/cantara/gober/stream"
	"github.com/cantara/gober/stream/event"
	"go/types"
)

type Store interface {
	Set(b store.Beer) (err error)
	Get(id string) (b store.Beer, err error)
}

type service struct {
	stream stream.Stream[store.Beer, types.Nil]
	store  Store
}

func InitService(st stream.Persistence, ctx context.Context) (s service, err error) {
	beerStream, err := stream.Init[store.Beer, types.Nil](st, "beer", ctx)
	if err != nil {
		panic(err)
	}
	beerStore, err := store.Init[store.Beer](beerStream, ctx)
	if err != nil {
		return
	}
	s = service{
		stream: beerStream,
		store:  beerStore,
	}
	return
}

func (s service) Get(id string) (b store.Beer, err error) {
	b, err = s.store.Get(id)
	log.AddError(err).Println(b, id)
	return
}

func (s service) Register(b store.Beer) (err error) {
	return s.store.Set(b)
}

func (s service) BeerStream(ctx context.Context) (out <-chan event.Event[store.Beer, types.Nil], err error) {
	out, err = s.stream.Stream(event.AllTypes(), streamStore.STREAM_START, stream.ReadAll[types.Nil](), prov, ctx)
	return
}
