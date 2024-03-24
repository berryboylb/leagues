package main

import (
	// "context"
	"flag"
	"log"
	"time"
	// "os"

	// apitoolkit "github.com/apitoolkit/apitoolkit-go"
	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/ratelimit"
	"github.com/fatih/color"
	"github.com/gin-contrib/cors"
	"github.com/joho/godotenv"
	"go.uber.org/ratelimit"
	"league/db"
)

var (
	limit = ratelimit.New(100)
	// rps   = flag.Int("rps", 100, "request per second")
)

func leakBucket() gin.HandlerFunc {
	prev := time.Now()
	return func(ctx *gin.Context) {
		now := limit.Take()
		log.Print(color.CyanString("%v", now.Sub(prev)))
		prev = now
	}
}

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	// ctx :=
	// apitoolkitClient, err := apitoolkit.NewClient(context.Background(), apitoolkit.Config{APIKey: os.Getenv("API_TOOLKIT")})
	// if err != nil {
	// 	log.Fatalf("Failed to load monitoring keys: %v", err)
	// }

	flag.Parse()

	app := gin.New()
	// app.Use(apitoolkitClient.GinMiddleware)
	app.Use(cors.Default())
	app.Use(leakBucket())
	router := app.Group("/api/v1")

	AddRoutes(router)

	// connect db
	db.ConnectDB()

	app.Run(":3000")

	log.Print("Server listening on http://localhost:3000/")
}
