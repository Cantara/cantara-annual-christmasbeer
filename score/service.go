package score

import (
	"context"

	"github.com/cantara/cantara-annual-christmasbeer/score/store"
	"github.com/cantara/gober/stream"
	"github.com/cantara/gober/stream/event"
	"github.com/gofrs/uuid"
)

type Store interface {
	Set(b store.Score) (err error)
	Get(id string) (b store.Score, err error)
	Stream(ctx context.Context) (out <-chan event.Event[store.Score], err error)
	Range(f func(key string, data store.Score) error)
}

type Account interface {
	Weight(id uuid.UUID) float32
}

type service struct {
	store   Store
	account Account
}

func InitService(st stream.Stream, a Account, ctx context.Context) (s service, err error) {
	scoreStore, err := store.Init[store.Score](st, ctx)
	if err != nil {
		return
	}
	s = service{
		store:   scoreStore,
		account: a,
	}
	return
}

func (s service) Get(id string) (b store.Score, err error) {
	b, err = s.store.Get(id)
	return
}

func (s service) Range(f func(key string, data store.Score) error) {
	s.store.Range(f)
}

func (s service) Register(b store.Score) (err error) {
	b.Weight = s.account.Weight(b.ScorerId)
	b.Newbie = b.Weight < 1
	b.Rating = float32(b.RatingBase) * b.Weight
	return s.store.Set(b)
}

func (s service) ScoreStream(ctx context.Context) (out <-chan event.Event[store.Score], err error) {
	//year := time.Now().Year()
	//md.Event.Year > year || md.Event.Year < year
	return s.store.Stream(ctx)
}
