package database

import (
	"context"
	"os"

	"github.com/redis/go-redis/v9"
)

var RDBClient *redis.Client

func RedisConn() {
	addr := os.Getenv("REDIS_HOST") + ":" + os.Getenv("REDIS_PORT")

	RDBClient = redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	_, err := RDBClient.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}
}

func RedisInstance() *redis.Client {
	if RDBClient == nil {
		RedisConn()
	}
	return RDBClient
}
