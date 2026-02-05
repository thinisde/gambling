package cache

import (
	"context"
	"os"
	"time"

	"github.com/go-redis/redis/v8"
)

var ctx = context.Background()

var CacheClient *RedisClient

type RedisClient struct {
	client *redis.Client
}

func NewRedisClient() *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	return &RedisClient{client: rdb}
}

func (r *RedisClient) Set(key string, value any, expiration time.Duration) error {
	return r.client.Set(ctx, key, value, expiration).Err()
}

func (r *RedisClient) Get(key string) (string, error) {
	return r.client.Get(ctx, key).Result()
}

func (r *RedisClient) Del(key string) error {
	return r.client.Del(ctx, key).Err()
}

func (r *RedisClient) Close() error {
	return r.client.Close()
}

func (r *RedisClient) Publish(channel string, message any) error {
	return r.client.Publish(ctx, channel, message).Err()
}

func (r *RedisClient) Subscribe(channel string) *redis.PubSub {
	return r.client.Subscribe(ctx, channel)
}

func (r *RedisClient) HSet(key string, values ...any) error {
	return r.client.HSet(ctx, key, values...).Err()
}

func (r *RedisClient) HGet(key, field string) (string, error) {
	return r.client.HGet(ctx, key, field).Result()
}

func (r *RedisClient) HDel(key string, fields ...string) error {
	return r.client.HDel(ctx, key, fields...).Err()
}

func (r *RedisClient) ClosePubSub(pubsub *redis.PubSub) error {
	return pubsub.Close()
}

func (r *RedisClient) LPush(key string, values ...any) error {
	return r.client.LPush(ctx, key, values...).Err()
}

func (r *RedisClient) RPop(key string) (string, error) {
	return r.client.RPop(ctx, key).Result()
}

func (r *RedisClient) LRange(key string, start, stop int64) ([]string, error) {
	return r.client.LRange(ctx, key, start, stop).Result()
}

func (r *RedisClient) LLen(key string) (int64, error) {
	return r.client.LLen(ctx, key).Result()
}

func (r *RedisClient) SAdd(key string, members ...any) error {
	return r.client.SAdd(ctx, key, members...).Err()
}

func (r *RedisClient) SRem(key string, members ...any) error {
	return r.client.SRem(ctx, key, members...).Err()
}
