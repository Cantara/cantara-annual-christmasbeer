package beer

type service struct {
}

func InitService() service {
	return service{}
}

func (s service) Get(id string) (b beer, err error) {
	return
}

func (s service) Register(b beer) (err error) {
	return
}
