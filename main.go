package main

import (
	"context"
	"github.com/cantara/cantara-annual-christmasbeer/account"
	"github.com/cantara/cantara-annual-christmasbeer/account/session"
	"github.com/cantara/cantara-annual-christmasbeer/account/store"
	"net"
	"net/http"
	"os"
	"time"

	log "github.com/cantara/bragi"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
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

func main() {
	loadEnv()
	since := time.Now()

	log.Println("Initialized webserver")
	r := gin.Default()
	config := cors.DefaultConfig()
	config.AllowOrigins = []string{"*"}
	r.Use(cors.New(config))
	base := r.Group("")

	dash := base.Group("/") //Might need to be in subdir dash
	{
		dash.StaticFile("/", "./frontend"+os.Getenv("frontend_path")+"/index.html")
		dash.StaticFile("/global.css", "./frontend"+os.Getenv("frontend_path")+"/global.css")
		dash.StaticFile("/favicon.png", "./frontend"+os.Getenv("frontend_path")+"/favicon.png")
		dash.StaticFS("/build", http.Dir("./frontend"+os.Getenv("frontend_path")+"/build"))
	}

	outboudIp := GetOutboundIP()
	api := base.Group("")
	api.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{
			"status":        "UP",
			"version":       Version,
			"build_time":    BuildTime,
			"ip":            outboudIp.String(),
			"running_since": since,
			"now":           time.Now(),
		})
	})

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
	accResource, err := account.InitResource(r, os.Getenv("api_path")+"/account", accService)
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

	r.Run(":" + os.Getenv("port"))
}
