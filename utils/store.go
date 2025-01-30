package utils

import (
	"context"
	"fmt"
	"os"

	"github.com/go-redis/redis/v8"
)

func NewRedisClient() *redis.Client {
	fmt.Println("Connecting to redis server: ", os.Getenv("REDIS_HOST"))

	// Create new redis client
	client := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return client
}

func SetKey(ctx *context.Context, client *redis.Client, key string, value string, ttl int) {
	fmt.Println("Setting the key", key, "to", value, "in Redis")

	// ttl = Time to live is set to 0, no expiration
	client.Set(*ctx, key, value, 0)

	fmt.Printf("Succesfully set key: %s -> value: %s\n", key, value)
}

func GetLongURL(ctx *context.Context, client *redis.Client, shortURL string) (string, error) {
	longURL, err := client.Get(*ctx, shortURL).Result()

	if err != redis.Nil {
		return "", fmt.Errorf("short URL not found")
	} else if err != nil {
		return "", fmt.Errorf("failed to retrieve url from redis server: %s", err)
	}

	return longURL, nil
}
