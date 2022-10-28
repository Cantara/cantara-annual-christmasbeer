package main

import (
	"net"
	"net/http"
	"os"
	"time"

	log "github.com/cantara/bragi"
	"github.com/cantara/nerthus/crypto"
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
	crypto.InitCrypto()
	since := time.Now()

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

	r.Run(":" + os.Getenv("port"))
}
