package main

import (
	"context"
	"fmt"
	log "github.com/cantara/bragi"
	"github.com/cantara/cantara-annual-christmasbeer/account"
	"github.com/cantara/cantara-annual-christmasbeer/account/session"
	"github.com/cantara/cantara-annual-christmasbeer/account/store"
	"github.com/cantara/gober"
	evStore "github.com/cantara/gober/store"
	"github.com/cantara/gober/store/inmemory"
	"github.com/cantara/gober/webserver"
	"github.com/cantara/gober/websocket"
	"github.com/joho/godotenv"
	"go/types"
	"net"
	"net/http"
	ws "nhooyr.io/websocket"
	"nhooyr.io/websocket/wsjson"
	"os"
	"time"
)

var Version string

var BuildTime string

func loadEnv() {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalf("Error loading .env file: %s", err)
	}
}

func GetOutboundIP() net.IP {
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	localAddr := conn.LocalAddr().(*net.UDPAddr)

	return localAddr.IP
}

type beer struct {
	Name  string `json:"name"`
	Brand string `json:"brand"`
	ABV   string `json:"abv"`
}

func prov(key string) string {
	return "MdgKIHmlbRszXjLbS7pXnSBdvl+SR1bSejtpFTQXxro="
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
	accSession, err := session.Init(ctx)
	if err != nil {
		panic(err)
	}
	log.Println("Initialized account session")
	accService, err := account.InitService(accStore, accSession, ctx)
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

	s, err := inmemory.Init()
	if err != nil {
		panic(err)
	}
	es, err := gober.Init[beer, types.Nil](s, "beer", ctx)
	if err != nil {
		panic(err)
	}
	go func() {
		i := 0
		for {
			es.Store(gober.Event[beer, types.Nil]{
				Type: "create",
				Data: beer{
					Name:  fmt.Sprintf("Test%d", i),
					Brand: "eXOReaction",
					ABV:   "5.8%",
				},
			}, prov)
			i++
			time.Sleep(10 * time.Second)
		}
	}()

	websocket.Websocket(api, "/beer", func(ctx context.Context, conn *ws.Conn) bool {
		conn.CloseRead(ctx)
		//ctxCancel, cancel := context.WithCancel(ctx)
		stream, err := es.Stream([]string{"create", "update", "delete"}, evStore.STREAM_START, gober.ReadAll[types.Nil](), prov, ctx)
		if err != nil {
			log.AddError(err).Error("while starting beer stream")
		}
		for e := range stream {
			err = wsjson.Write(ctx, conn, e.Data)
			if err != nil {
				log.AddError(err).Warning("while writing to socket")
				return false
			}
		}
		return false
	})

	serv.Run()
}
