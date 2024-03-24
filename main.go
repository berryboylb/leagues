package main

import (
	// "context"
	"log"
	"flag"
	"time"
	// "os"

	// apitoolkit "github.com/apitoolkit/apitoolkit-go"
	"github.com/gin-gonic/gin"
	// "github.com/gin-contrib/ratelimit"
    "github.com/gin-contrib/cors" 
	"github.com/joho/godotenv"
	"league/db"
	"github.com/fatih/color"
	"go.uber.org/ratelimit"
)

var (
	limit ratelimit.Limiter
	rps   = flag.Int("rps", 100, "request per second")
)

func init() {
	log.SetFlags(0)
	log.SetPrefix("[GIN] ")
	log.SetOutput(gin.DefaultWriter)
}

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

	limit = ratelimit.New(*rps)

	app := gin.New()
	// app.Use(apitoolkitClient.GinMiddleware)
	app.Use(cors.Default())
	app.Use(leakBucket())
	router := app.Group("/api/v1")


    // // Apply rate limiting middleware
    // router.Use(ratelimit.New(ratelimit.IPRateLimiter(10, 1*time.Minute)))

	AddRoutes(router)

	log.Printf(color.CyanString("Current Rate Limit: %v requests/s", *rps))

	// connect db
	db.ConnectDB()

	app.Run(":3000")

	log.Print("Server listening on http://localhost:3000/")
}
