package database

import (
	"context"
	"log"
	"time"

	"github.com/go-redis/redis/v8"
)

var RedisClient *redis.Client
var Ctx = context.Background()

func InitializeRedis() {
	RedisClient = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",               
		DB:       0,                
	})

	_, err := RedisClient.Ping(Ctx).Result()
	if err != nil {
		log.Fatalf("Failed to connect to Redis: %v", err)
	} else {
		log.Println("Connected to Redis successfully!")
	}
}

func SetCache(key string, value string, ttl time.Duration) error {
	return RedisClient.Set(Ctx, key, value, ttl).Err()
}

func GetCache(key string) (string, error) {
	return RedisClient.Get(Ctx, key).Result()
}

func DeleteCache(key string) error {
	return RedisClient.Del(Ctx, key).Err()
}
