package db

import (
	"github.com/joho/godotenv"

	"context"
	"fmt"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
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

	ctx, _ := context.WithTimeout(context.Background(), 30*time.Second)
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

func IndexField(collection mongo.Collection, field string, indexType int) error {
	indexModel := mongo.IndexModel{
		Keys:    bson.M{field: indexType}, // 1 for ascending, -1 for descending
		Options: options.Index().SetUnique(true),
	}
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	_, err := collection.Indexes().CreateOne(ctx, indexModel)
	if err != nil {
		fmt.Println("Failed to create index:", err)
		return err
	}
	fmt.Println("Index created successfully on %v field.", field)
	return nil
}

func IsIndexExists(ctx context.Context, collection *mongo.Collection, indexKey string) (bool, error) {
    indexes, err := collection.Indexes().List(ctx)
    if err != nil {
        return false, err
    }
    defer indexes.Close(ctx)

    for indexes.Next(ctx) {
        var index bson.M
        if err := indexes.Decode(&index); err != nil {
            return false, err
        }
        if val, ok := index["key"].(bson.M)[indexKey]; ok && val != nil {
            return true, nil
        }
    }
    return false, nil
}