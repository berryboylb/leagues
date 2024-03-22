package db

import (
	"github.com/joho/godotenv"

	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)


func EnvMongoURI() string {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }
	uri := os.Getenv("MONGO_URI")
	if uri == "" {
		log.Fatal("You must set your 'MONGODB_URI' environment variable.")
	}
    return uri
}

func ConnectDB() *mongo.Client {
	// fmt.Println(uri, "nnhh")
	client, err := mongo.NewClient(options.Client().ApplyURI(EnvMongoURI()))
	if err != nil {
		log.Fatal(err)
	}

	ctx, _ := context.WithTimeout(context.Background(), 10*time.Second)
	err = client.Connect(ctx)
	if err != nil {
		log.Fatal(err)
	}

	//ping the database
	err = client.Ping(ctx, nil)
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println("Connected to MongoDB")
	return client
}

// Client instance
var MongoClient *mongo.Client = ConnectDB()

// getting database collections
func GetCollection(client *mongo.Client, collectionName string) *mongo.Collection {
	// the database name on atlas
	collection := client.Database("league").Collection(collectionName)
	return collection
}
