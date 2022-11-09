package main

import (
	"context"
	log "github.com/cantara/bragi"
	"github.com/cantara/cantara-annual-christmasbeer/account"
	"github.com/cantara/cantara-annual-christmasbeer/account/admin"
	"github.com/cantara/cantara-annual-christmasbeer/account/session"
	"github.com/cantara/cantara-annual-christmasbeer/account/store"
	"github.com/cantara/cantara-annual-christmasbeer/beer"
	"github.com/cantara/gober/webserver"
	"github.com/joho/godotenv"
	"net/http"
	"os"
)

var Version string

var BuildTime string

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func main() {
	loadEnv()

	serv := webserver.Init()
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
	accStore, err := store.Init(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Initialized account store")
	admStore, err := admin.Init[struct{}](ctx)
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

	beerService := beer.InitService()
	_, err = beer.InitResource(api, "/beer", accService, beerService, ctx)

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
		err = accService.RegisterAdmin(accountId)
		if err != nil {
			panic(err)
		}
	}

	serv.Run()
}
