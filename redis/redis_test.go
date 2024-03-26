package redis 

import (
	"fmt"
	"testing"
	"time"
	"github.com/stretchr/testify/assert"
	"errors"
	"github.com/go-redis/redis/v8"
)



func setupRedisTestEnvironment() {
	// Initialize Redis client for testing
	client = redis.NewClient(&redis.Options{
		Addr: "localhost:6379", // Assuming Redis is running on localhost
		Password: "",
		DB:       0,
	})

	// Ping Redis to ensure connection
	_, err := client.Ping(ctx).Result()
	if err != nil {
		panic(fmt.Sprintf("failed to connect to Redis: %v", err))
	}
}

func cleanupRedisTestEnvironment() {
	// Close Redis client after testing
	if client != nil {
		if err := client.Close(); err != nil {
			fmt.Printf("failed to close Redis client: %v\n", err)
		}
	}
}


func TestRedisFunctions(t *testing.T) {
	// Set up test environment
	setupRedisTestEnvironment()
	defer cleanupRedisTestEnvironment()

	// Test data
	key := "test_key"
	value := "test_value"
	expiration := time.Hour

	// Store value in Redis
	err := Store(key, value, expiration)
	assert.NoError(t, err)

	// Retrieve value from Redis
	retrievedValue, err := Retrieve(key)
	assert.NoError(t, err)
	assert.Equal(t, value, retrievedValue)

	// Delete value from Redis
	err = Delete(key)
	assert.NoError(t, err)

	// Attempt to retrieve deleted value
	_, err = Retrieve(key)
	assert.Error(t, err)
	assert.True(t, errors.Is(err, redis.Nil))
}
