package redis

import (
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"time"
	"encoding/json"
)

var client *redis.Client
var ctx = context.Background()


func init() {
	err := godotenv.Load()
	if err != nil {
		log.Println("Error loading .env file")
	}

	redisPassword := os.Getenv("REDIS_PASSWORD")
	redisDB := os.Getenv("REDIS_DB")
	redisHost := os.Getenv("REDIS_HOST")
	fmt.Println(redisPassword, redisDB, redisHost)

	if  redisDB == "" || redisHost == "" {
		panic("Error loading redis password, db or host")
	}

	hostNumber, err := strconv.Atoi(redisDB)
	if err != nil {
		panic("Error converting redis host to int")
	}
	client = redis.NewClient(&redis.Options{
		Addr:     redisHost,
		Password: redisPassword, 
		DB:       hostNumber,    
	})

	Test()
}
// exposes client to be used by other packages
func GetClient() *redis.Client {
	return client
}


// Test function to test connection to redis
func Test() {
	ping, err := client.Ping(ctx).Result()
	if err != nil {
		fmt.Println(err.Error())
		panic(err)
		return
	}
	fmt.Println(ping)
}

// Store function to store a value in redis
func Store(key string, value interface{}, expiration time.Duration) error {
	err := client.Set(ctx, key, value, expiration).Err()
	if err != nil {
		fmt.Printf("failed to set value in redis instance: %v\n", err)
		return err
	}
	return nil
}

// StoreStruct serializes a struct into a JSON string.
func StoreStruct(value interface{}) ([]byte, error) {
	jsonValue, err := json.Marshal(value)
	if err != nil {
		return nil, err
	}
	return jsonValue, nil
}


// UnmarshalStruct unserializes a JSON string into a struct.
func UnmarshalStruct(jsonValue []byte, result interface{}) error {
	err := json.Unmarshal(jsonValue, &result)
	if err != nil {
		return err
	}
	return nil
}


func Retrieve(key string) (string, error) {
	result, err := client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			fmt.Printf("key %s not found in Redis\n", key)
			return "", err
		}
		fmt.Printf("failed to get value for key %v from Redis: %v", key, err)
		return "", err
	}
	return result, nil
}
