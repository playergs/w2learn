package database

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
)

type RedisConfig struct {
	Addr         string
	Password     string
	DB           int
	PoolSize     int
	MinIdleConns int
	MaxRetries   int
	DialTimeout  time.Duration
	ReadTimeout  time.Duration
	WriteTimeout time.Duration
}

var redisClient *redis.Client

func NewRedis(config *RedisConfig) (*redis.Client, error) {
	if redisClient != nil {
		return redisClient, nil
	}

	client := redis.NewClient(&redis.Options{
		Addr:         config.Addr,
		Password:     config.Password,
		DB:           config.DB,
		PoolSize:     config.PoolSize,
		MinIdleConns: config.MinIdleConns,
		MaxRetries:   config.MaxRetries,
		DialTimeout:  config.DialTimeout,
		ReadTimeout:  config.ReadTimeout,
		WriteTimeout: config.WriteTimeout,
	})

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)

	defer cancel()

	err := client.Ping(ctx).Err()

	if err != nil {
		return nil, err
	}

	redisClient = client

	return redisClient, nil
}

func GetRedis() *redis.Client {
	return redisClient
}

func CloseRedis() error {
	if redisClient == nil {
		return nil
	}

	return redisClient.Close()
}
