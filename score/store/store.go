package store

import (
	"context"
	"fmt"
	"github.com/cantara/cantara-annual-christmasbeer/beer/store"
	"github.com/cantara/gober/persistenteventmap"
	"github.com/cantara/gober/stream"
	"github.com/gofrs/uuid"
)

type storeService[pt any] struct {
	scores persistenteventmap.EventMap[Score, ScoreMetadata]
}

func CryptoKey(_ string) string {
	return "a1BNjgicHSQ/YKgG8qhvi9I2MdZXcXQNoef7bnimDLI="
}

func Init[pt any](st stream.Stream[Score, ScoreMetadata], ctx context.Context) (s storeService[pt], err error) {
	scoreMap, err := persistenteventmap.Init[Score, ScoreMetadata](st, "score", "0.1.0", CryptoKey, func(s Score) string {
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
	err = s.scores.Set(b, ScoreMetadata{
		Year: b.Year,
	})
	if err != nil {
		return
	}
	return
}

func (s storeService[pt]) Get(id string) (b Score, err error) {
	b, _, err = s.scores.Get(id)
	return
}

type Score struct {
	Scorer  uuid.UUID  `json:"scorer"`
	Year    int        `json:"year"`
	Beer    store.Beer `json:"beer"`
	Rating  int        `json:"rating"`
	Comment string     `json:"comment"`
}

func (s Score) ToId() string {
	return fmt.Sprintf("%s_%d_%s", s.Scorer, s.Year, s.Beer.ToId())
}

type ScoreMetadata struct {
	Year int `json:"year"`
}
