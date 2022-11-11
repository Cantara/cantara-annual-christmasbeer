package score

import (
	"context"
	log "github.com/cantara/bragi"
	"github.com/cantara/cantara-annual-christmasbeer/score/store"
	streamStore "github.com/cantara/gober/store"
	"github.com/cantara/gober/stream"
	"github.com/cantara/gober/stream/event"
	"time"
)

type Store interface {
	Set(b store.Score) (err error)
	Get(id string) (b store.Score, err error)
}

type service struct {
	stream stream.Stream[store.Score, store.ScoreMetadata]
	store  Store
}

func InitService(st stream.Persistence, ctx context.Context) (s service, err error) {
	scoreStream, err := stream.Init[store.Score, store.ScoreMetadata](st, "score", ctx)
	if err != nil {
		panic(err)
	}
	scoreStore, err := store.Init[store.Score](scoreStream, ctx)
	if err != nil {
		return
	}
	s = service{
		stream: scoreStream,
		store:  scoreStore,
	}
	return
}

func (s service) Get(id string) (b store.Score, err error) {
	b, err = s.store.Get(id)
	return
}

func (s service) Register(b store.Score) (err error) {
	return s.store.Set(b)
}

func (s service) BeerStream(ctx context.Context) (out <-chan event.Event[store.Score, store.ScoreMetadata], err error) {
	year := time.Now().Year()
	out, err = s.stream.Stream(event.AllTypes(), streamStore.STREAM_START, func(md event.Metadata[store.ScoreMetadata]) bool {
		log.Debug(md)
		if md.DataType != "score" {
			return true
		}
		if md.Event.Year > year || md.Event.Year < year {
			return true
		}
		return false
	}, store.CryptoKey, ctx)
	return
}
