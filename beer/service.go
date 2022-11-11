package beer

import "github.com/cantara/cantara-annual-christmasbeer/beer/store"

type service struct {
}

func InitService() service {
	return service{}
}

func (s service) Get(id string) (b store.beer, err error) {
	return
}

func (s service) Register(b store.beer) (err error) {
	return
}
