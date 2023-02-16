package store

import (
	"context"
	"fmt"
	"github.com/cantara/cantara-annual-christmasbeer/beer/store"
	"github.com/cantara/gober/persistenteventmap"
	"github.com/cantara/gober/stream"
	"github.com/cantara/gober/stream/event"
	eventStore "github.com/cantara/gober/stream/event/store"
	"github.com/gofrs/uuid"
)

type storeService[pt any] struct {
	scores persistenteventmap.EventMap[Score]
}

func CryptoKey(_ string) string {
	return "a1BNjgicHSQ/YKgG8qhvi9I2MdZXcXQNoef7bnimDLI="
}

func Init[pt any](st stream.Stream, ctx context.Context) (s storeService[pt], err error) {
	scoreMap, err := persistenteventmap.Init[Score](st, "score", "0.1.0", CryptoKey, func(s Score) string {
		return s.ToId()
	}, ctx)
	if err != nil {
		return
	}
	s = storeService[pt]{
		scores: scoreMap,
	}
	return
}

func (s storeService[pt]) Set(b Score) (err error) {
	err = s.scores.Set(b)
	if err != nil {
		return
	}
	return
}

func (s storeService[pt]) Get(id string) (b Score, err error) {
	b, err = s.scores.Get(id)
	return
}

func (s storeService[pt]) Stream(ctx context.Context) (out <-chan event.Event[Score], err error) {
	return s.scores.Stream(event.AllTypes(), eventStore.STREAM_START, stream.ReadDataType("score"), ctx)
}

type Score struct {
	ScorerId   uuid.UUID  `json:"scorer_id"`
	Scorer     string     `json:"scorer"`
	Year       int        `json:"year"`
	Beer       store.Beer `json:"beer"`
	Rating     float32    `json:"rating"`
	RatingBase int        `json:"rating_base"`
	Weight     float32    `json:"weight"`
	Newbie     bool       `json:"newbie"`
	Comment    string     `json:"comment"`
}

func (s Score) ToId() string {
	return fmt.Sprintf("%s_%d_%s", s.ScorerId, s.Year, s.Beer.ToId())
}
