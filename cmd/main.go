package main

import (
	"encoding/json"
	"fmt"
	"go/types"
	"io"
	"os"
	"strings"

	log "github.com/cantara/bragi"
	"github.com/cantara/cantara-annual-christmasbeer/account"
	bs "github.com/cantara/cantara-annual-christmasbeer/beer/store"
	ss "github.com/cantara/cantara-annual-christmasbeer/score/store"
	"github.com/cantara/gober/stream/event"
	"github.com/dgraph-io/badger"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func storeToCSV[T any](store string) (map[string]T, error) {
	db, err := badger.Open(badger.DefaultOptions("./eventmap/" + store))
	defer db.Close()
	if err != nil {
		return nil, err
	}
	f, err := os.OpenFile(store+".csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	io.WriteString(f, "key,value\n")
	o := make(map[string]T)
	err = db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.PrefetchSize = 10
		it := txn.NewIterator(opts)
		defer it.Close()
		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			k := item.Key()
			if strings.HasSuffix(string(k), "position") {
				return nil
			}
			err := item.Value(func(v []byte) error {
				var t T
				json.Unmarshal(v, &t)
				o[string(k)] = t
				fmt.Fprintf(f, "%s,%s\n", k, v)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	return o, err
}

type dmd[DT, MT any] struct {
	Data     DT                 `json:"data"`
	Metadata event.Metadata[MT] `json:"metadata"`
}

func main() {
	loadEnv()
	log.SetLevel(log.INFO)
	beers, _ := storeToCSV[dmd[bs.Beer, types.Nil]]("beer")
	scores, _ := storeToCSV[dmd[ss.Score, ss.ScoreMetadata]]("score")
	accs, _ := storeToCSV[dmd[account.AccountRegister, types.Nil]]("account")
	f, err := os.OpenFile("result.csv", os.O_TRUNC|os.O_CREATE|os.O_WRONLY, 0640)
	if err != nil {
		return
	}
	fmt.Fprint(f, ",")
	for _, beer := range beers {
		fmt.Fprintf(f, "%s,", beer.Data.Name)
		fmt.Println("bn", beer.Data.Name)
	}
	fmt.Fprintln(f)
	for id, acc := range accs {
		if acc.Data.FirstName == "" {
			continue
		}
		fmt.Fprintf(f, "%s,", acc.Data.FirstName)
	beer:
		for _, beer := range beers {
			for _, score := range scores {
				if beer.Data.Name != score.Data.Beer.Name {
					continue
				}
				if score.Data.ScorerId.String() == id {
					fmt.Println("bn", beer.Data.Name, "s", int(score.Data.Rating))
					fmt.Fprintf(f, "%d,", int(score.Data.Rating))
					continue beer
				} else {
					continue
				}
			}
			fmt.Fprint(f, ",")
		}
		fmt.Fprintln(f)
	}
}
