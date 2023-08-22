package core

import (
	"context"
	"encoding/json"
	"time"

	"github.com/go-redis/redis/v8"
	"github.com/spf13/viper"
)

type RedisClient interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string) ([]byte, error)
	LPush(ctx context.Context, key string, values ...interface{}) error
	GetList(ctx context.Context, key string) ([]string, error)
}

type redisClient struct {
	rdb *redis.Client
}

// Set menyimpan data ke Redis dengan key dan value yang diberikan serta durasi expiration
func (c *redisClient) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	jsonData, err := json.Marshal(value)
	if err != nil {
		return err
	}

	err = c.rdb.Set(ctx, key, jsonData, expiration).Err()
	return err
}

// Get mengambil data dari Redis berdasarkan key yang diberikan
func (c *redisClient) Get(ctx context.Context, key string) ([]byte, error) {
	val, err := c.rdb.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	return val, nil
}
func (c *redisClient) GetList(ctx context.Context, key string) ([]string, error) {
	val, err := c.rdb.LRange(ctx, key, 0, -1).Result()
	if err != nil {
		return nil, err
	}

	return val, nil
}

// LPush menambahkan satu atau lebih value ke awal sebuah list dengan key yang diberikan
func (c *redisClient) LPush(ctx context.Context, key string, values ...interface{}) error {
	err := c.rdb.LPush(ctx, key, values...).Err()
	return err
}

// NewRedisClient membuat instance RedisClient menggunakan konfigurasi yang ada di environment variable
func NewRedisClient() RedisClient {
	client := redis.NewClient(&redis.Options{
		Addr:     viper.GetString("redis_credentials.REDIS_URL"),      // alamat Redis
		Password: viper.GetString("redis_credentials.REDIS_PASSWORD"), // password Redis (kosong jika tidak ada)
		DB:       0,                                                   // nomor database Redis
	})

	_, err := client.Ping(context.Background()).Result()
	if err != nil {
		panic(err)
	}

	return &redisClient{rdb: client}
}
