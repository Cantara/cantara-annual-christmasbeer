package main

import (
	"context"
	"embed"
	"fmt"
	"io"
	"io/fs"
	"net/http"
	"os"
	"strconv"
	"time"

	log "github.com/cantara/bragi"
	"github.com/cantara/bragi/sbragi"
	"github.com/cantara/cantara-annual-christmasbeer/account"
	"github.com/cantara/cantara-annual-christmasbeer/account/privilege"
	"github.com/cantara/cantara-annual-christmasbeer/account/session"
	"github.com/cantara/cantara-annual-christmasbeer/account/store"
	"github.com/cantara/cantara-annual-christmasbeer/beer"
	"github.com/cantara/cantara-annual-christmasbeer/score"
	"github.com/cantara/gober/stream"
	"github.com/cantara/gober/stream/event/store/eventstore"
	"github.com/cantara/gober/stream/event/store/ondisk"
	"github.com/cantara/gober/webserver"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/gofrs/uuid"
	"github.com/joho/godotenv"
)

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

//go:embed frontend/static/*
var static embed.FS

//go:embed frontend/*.html
var pages embed.FS

func main() {
	loadEnv()
	log.SetLevel(log.INFO)
	pages, err := fs.Sub(pages, "frontend")
	if err != nil {
		sbragi.WithError(err).Fatal("while creating sub of embed fs")
	}
	fs.WalkDir(pages, ".", func(path string, d fs.DirEntry, err error) error {
		log.Info(path)
		return nil
	})
	static, err := fs.Sub(static, "frontend/static")
	if err != nil {
		sbragi.WithError(err).Fatal("while creating sub of embed fs")
	}
	fs.WalkDir(static, ".", func(path string, d fs.DirEntry, err error) error {
		log.Info(path)
		return nil
	})

	portString := os.Getenv("webserver.port")
	port, err := strconv.Atoi(portString)
	if err != nil {
		log.AddError(err).Fatal("while getting webserver port")
	}
	serv, err := webserver.Init(uint16(port), false)
	if err != nil {
		log.AddError(err).Fatal("while initializing webserver")
	}
	log.Println("Initialized webserver")

	api := serv.API()
	api.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3030"},
		AllowMethods:     []string{"PUT", "PATCH"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

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
		accountStream, err = ondisk.Init("account", ctx)
		if err != nil {
			log.AddError(err).Fatal("while creating account stream")
		}
		adminStream, err = ondisk.Init("admin", ctx)
		if err != nil {
			log.AddError(err).Fatal("while creating admin stream")
		}
		beerStream, err = ondisk.Init("beer", ctx)
		if err != nil {
			log.AddError(err).Fatal("while creating beer stream")
		}
		scoreStream, err = ondisk.Init("score", ctx)
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

	//Adding frontend content
	{
		indexF, err := pages.Open("index.html")
		if err != nil {
			sbragi.WithError(err).Fatal("while opeining index")
		}
		indexB, err := io.ReadAll(indexF)
		if err != nil {
			indexF.Close()
			sbragi.WithError(err).Fatal("while reading index")
		}
		index := string(indexB)
		indexF.Close()
		votingF, err := pages.Open("voting.html")
		if err != nil {
			sbragi.WithError(err).Fatal("while opeining voting")
		}
		votingB, err := io.ReadAll(votingF)
		if err != nil {
			indexF.Close()
			sbragi.WithError(err).Fatal("while reading voting")
		}
		voting := string(votingB)
		votingF.Close()
		api.GET("/", func(c *gin.Context) {
			if accResource.Account(c) != uuid.Nil {
				fmt.Fprint(c.Writer, voting)
				return
			}
			fmt.Fprint(c.Writer, index)
			c.Writer.WriteHeader(http.StatusOK)
		})
		api.StaticFS("/static", http.FS(static))
		//api.StaticFile("/", "./frontend"+os.Getenv("frontend_path")+"/index.html")
		//api.StaticFile("/global.css", "./frontend"+os.Getenv("frontend_path")+"/global.css")
		//api.StaticFile("/favicon.png", "./frontend"+os.Getenv("frontend_path")+"/favicon.png")
		//api.StaticFS("/build", http.Dir("./frontend"+os.Getenv("frontend_path")+"/build"))
	}

	serv.Run()
}
