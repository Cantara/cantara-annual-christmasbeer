package main

import (
	"context"
	log "github.com/cantara/bragi"
	"github.com/cantara/cantara-annual-christmasbeer/account"
	"github.com/cantara/cantara-annual-christmasbeer/account/privilege"
	"github.com/cantara/cantara-annual-christmasbeer/account/session"
	"github.com/cantara/cantara-annual-christmasbeer/account/store"
	"github.com/cantara/cantara-annual-christmasbeer/beer"
	"github.com/cantara/cantara-annual-christmasbeer/score"
	"github.com/cantara/gober/stream"
	"github.com/cantara/gober/stream/event/store/eventstore"
	"github.com/cantara/gober/stream/event/store/inmemory"
	"github.com/cantara/gober/webserver"
	"github.com/joho/godotenv"
	"net/http"
	"os"
	"strconv"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func main() {
	loadEnv()
	log.SetLevel(log.INFO)

	portString := os.Getenv("webserver.port")
	port, err := strconv.Atoi(portString)
	if err != nil {
		log.AddError(err).Fatal("while getting webserver port")
	}
	serv, err := webserver.Init(uint16(port))
	if err != nil {
		log.AddError(err).Fatal("while initializing webserver")
	}
	log.Println("Initialized webserver")

	api := serv.API
	{
		api.StaticFile("/", "./frontend"+os.Getenv("frontend_path")+"/index.html")
		api.StaticFile("/global.css", "./frontend"+os.Getenv("frontend_path")+"/global.css")
		api.StaticFile("/favicon.png", "./frontend"+os.Getenv("frontend_path")+"/favicon.png")
		api.StaticFS("/build", http.Dir("./frontend"+os.Getenv("frontend_path")+"/build"))
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()
	var accountStream stream.Stream
	var adminStream stream.Stream
	var beerStream stream.Stream
	var scoreStream stream.Stream
	if esHost := os.Getenv("eventstore.host"); len(esHost) > 1 {
		es, err := eventstore.NewClient(esHost)
		if err != nil {
			log.AddError(err).Fatal("while creating eventstore client")
		}
		accountStream, err = eventstore.NewStream(es, "account", ctx)
		if err != nil {
			log.AddError(err).Fatal("while creating account stream")
		}
		adminStream, err = eventstore.NewStream(es, "admin", ctx)
		if err != nil {
			log.AddError(err).Fatal("while creating admin stream")
		}
		beerStream, err = eventstore.NewStream(es, "beer", ctx)
		if err != nil {
			log.AddError(err).Fatal("while creating beer stream")
		}
		scoreStream, err = eventstore.NewStream(es, "score", ctx)
		if err != nil {
			log.AddError(err).Fatal("while creating score stream")
		}
	} else {
		var err error
		accountStream, err = inmemory.Init("account", ctx)
		if err != nil {
			log.AddError(err).Fatal("while creating account stream")
		}
		adminStream, err = inmemory.Init("admin", ctx)
		if err != nil {
			log.AddError(err).Fatal("while creating admin stream")
		}
		beerStream, err = inmemory.Init("beer", ctx)
		if err != nil {
			log.AddError(err).Fatal("while creating beer stream")
		}
		scoreStream, err = inmemory.Init("score", ctx)
		if err != nil {
			log.AddError(err).Fatal("while creating score stream")
		}
	}
	accStore, err := store.Init(accountStream, ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Initialized account store")
	admStore, err := privilege.Init[account.Rights](adminStream, ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Initialized admin store")
	accSession, err := session.Init(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Initialized account session")
	accService, err := account.InitService(accStore, admStore, accSession, ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Initialized account resource")
	accResource, err := account.InitResource(api, os.Getenv("api_path")+"/account", accService)
	if err != nil {
		panic(err)
	}
	log.Println("Initialized account endpoints")
	if os.Getenv("account.internal.enable") == "true" {
		err = accResource.InitResourceInternal()
		if err != nil {
			panic(err)
		}
		log.Println("Initialized account internal endpoints")
	}

	beerService, err := beer.InitService(beerStream, ctx)
	if err != nil {
		panic(err)
	}
	_, err = beer.InitResource(api, "/beer", accService, beerService, ctx)

	scoreService, err := score.InitService(scoreStream, accService, ctx)
	if err != nil {
		panic(err)
	}
	_, err = score.InitResource(api, "/score", accService, beerService, scoreService, ctx)

	log.Println("Checking if admin user exists")
	_, err = accService.GetByUsername(os.Getenv("admin.username"))
	if err == nil {
		log.Println("Admin user already exists")
	} else {
		log.Println("Registering predefined admin user")
		token, err := accService.Register(account.AccountRegister{
			Username:  os.Getenv("admin.username"),
			Email:     os.Getenv("admin.email"),
			FirstName: os.Getenv("admin.first_name"),
			LastName:  os.Getenv("admin.last_name"),
			Number:    os.Getenv("admin.number"),
			Password:  os.Getenv("admin.password"),
		})
		if err != nil {
			panic(err)
		}

		_, accountId, err := accService.Validate(token.Token)
		if err != nil {
			panic(err)
		}
		log.Println("Registering predefined admin rights")
		err = accService.RegisterAdmin(accountId, account.Rights{Admin: true})
		if err != nil {
			panic(err)
		}
	}

	serv.Run()
}
