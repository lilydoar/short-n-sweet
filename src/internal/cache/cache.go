package cache

import (
	"context"
	"time"

	"github.com/redis/go-redis/v9"
	log "github.com/sirupsen/logrus"
)

type CacheConfig struct {
	Address string `json:"address" yaml:"address" toml:"address" env:"CACHE_ADDRESS"`
}

const (
	cacheDuration = 1 * time.Hour
)

type Cache interface {
	Get(shortUrl string) (string, error)
	Set(shortUrl string, longUrl string) error
}

type RedisCache struct {
	client *redis.Client
}

func InitRedisCache(config CacheConfig) *RedisCache {
	redisCache := &RedisCache{}

	redisClient := redis.NewClient(&redis.Options{
		Addr:     config.Address,
		Password: "",
		DB:       0,
	})

	_, err := redisClient.Ping(context.Background()).Result()
	if err != nil {
		log.Fatal(err)
	} else {
		log.Info("Connected to Redis at " + config.Address)
	}

	redisCache.client = redisClient

	return redisCache
}

func (redisCache *RedisCache) Get(shortUrl string) (string, error) {
	longUrl, err := redisCache.client.Get(context.Background(), shortUrl).Result()
	if err != nil {
		return "", err
	}

	return longUrl, nil
}

func (redisCache *RedisCache) Set(shortUrl string, longUrl string) error {
	err := redisCache.client.Set(context.Background(), shortUrl, longUrl, cacheDuration).Err()
	if err != nil {
		return err
	}

	return nil
}
