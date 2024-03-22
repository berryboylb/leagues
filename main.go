package main

import (
	// "context"
	"log"
	// "os"

	// apitoolkit "github.com/apitoolkit/apitoolkit-go"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"league/db"
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatalf("Failed to load the env vars: %v", err)
	}

	// ctx := 
	// apitoolkitClient, err := apitoolkit.NewClient(context.Background(), apitoolkit.Config{APIKey: os.Getenv("API_TOOLKIT")})
	// if err != nil {
	// 	log.Fatalf("Failed to load monitoring keys: %v", err)
	// }

	app := gin.New()
	// app.Use(apitoolkitClient.GinMiddleware)
	router := app.Group("/api/v1")

	AddRoutes(router)

	// connect db
	db.ConnectDB()

	app.Run(":3000")

	log.Print("Server listening on http://localhost:3000/")
}
