package storage

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
)

// CacheService определяет интерфейс для работы с кешем
type CacheService interface {
	Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error
	Get(ctx context.Context, key string, dest interface{}) error
	Delete(ctx context.Context, key string) error
	DeletePattern(ctx context.Context, pattern string) error
}

// RedisCache реализует CacheService с использованием Redis
type RedisCache struct {
	client *redis.Client
}

// NewRedisCache создает новый экземпляр RedisCache
func NewRedisCache(redisURL string) (*RedisCache, error) {
	opt, err := redis.ParseURL(redisURL)
	if err != nil {
		return nil, fmt.Errorf("failed to parse Redis URL: %w", err)
	}

	client := redis.NewClient(opt)

	// Проверяем подключение
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := client.Ping(ctx).Err(); err != nil {
		return nil, fmt.Errorf("failed to connect to Redis: %w", err)
	}

	log.Println("Redis connection successful")
	return &RedisCache{client: client}, nil
}

// Set сохраняет значение в кеше
func (r *RedisCache) Set(ctx context.Context, key string, value interface{}, expiration time.Duration) error {
	data, err := json.Marshal(value)
	if err != nil {
		return fmt.Errorf("failed to marshal value: %w", err)
	}

	err = r.client.Set(ctx, key, data, expiration).Err()
	if err != nil {
		log.Printf("Error setting cache key %s: %v", key, err)
		return err
	}

	log.Printf("Cache set: %s (expires in %v)", key, expiration)
	return nil
}

// Get получает значение из кеша
func (r *RedisCache) Get(ctx context.Context, key string, dest interface{}) error {
	data, err := r.client.Get(ctx, key).Result()
	if err != nil {
		if err == redis.Nil {
			log.Printf("Cache miss: %s", key)
			return fmt.Errorf("key not found")
		}
		log.Printf("Error getting cache key %s: %v", key, err)
		return err
	}

	err = json.Unmarshal([]byte(data), dest)
	if err != nil {
		log.Printf("Error unmarshaling cache data for key %s: %v", key, err)
		return fmt.Errorf("failed to unmarshal cached data: %w", err)
	}

	log.Printf("Cache hit: %s", key)
	return nil
}

// Delete удаляет ключ из кеша
func (r *RedisCache) Delete(ctx context.Context, key string) error {
	err := r.client.Del(ctx, key).Err()
	if err != nil {
		log.Printf("Error deleting cache key %s: %v", key, err)
		return err
	}

	log.Printf("Cache deleted: %s", key)
	return nil
}

// DeletePattern удаляет все ключи по паттерну
func (r *RedisCache) DeletePattern(ctx context.Context, pattern string) error {
	keys, err := r.client.Keys(ctx, pattern).Result()
	if err != nil {
		log.Printf("Error finding keys with pattern %s: %v", pattern, err)
		return err
	}

	if len(keys) > 0 {
		err = r.client.Del(ctx, keys...).Err()
		if err != nil {
			log.Printf("Error deleting keys with pattern %s: %v", pattern, err)
			return err
		}
		log.Printf("Cache deleted %d keys with pattern: %s", len(keys), pattern)
	}

	return nil
}

// Close закрывает соединение с Redis
func (r *RedisCache) Close() error {
	return r.client.Close()
} 