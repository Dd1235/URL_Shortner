package store

import (
	"context"
	"fmt"
	"os"

	"github.com/redis/go-redis/v9"
)

type RedisStore struct {
	client *redis.Client
}

func NewRedisStore() *RedisStore {
	fmt.Println("Connecting to Redis:", os.Getenv("REDIS_HOST"))

	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_HOST"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0, // use default db
	})

	return &RedisStore{client: rdb} // return pointer to object which is of type struct RedisStore (sorry need to do some c style processing in my head, uff reminds me of telugu kannada)
}

// pointer receiver , so method for RedisStore
func (r *RedisStore) Set(ctx context.Context, key string, value string) error {
	fmt.Printf("Setting key %s to value %s\n", key, value)

	err := r.client.Set(ctx, key, value, 0).Err() // yes, pointer.val works for pointers too, automatically deref
	// 0 - there is no expiration time in the database for this entry
	fmt.Println("The key", key, "has been set to", value, "successfully")
	return err
}

func (r *RedisStore) Get(ctx context.Context, key string) (string, error) {
	val, err := r.client.Get(ctx, key).Result()
	if err == redis.Nil {
		return "", fmt.Errorf("key not found")
	} else if err != nil {
		return "", err
	}
	return val, nil
}

func (r *RedisStore) Close() error {
	return r.client.Close()
}

func (r *RedisStore) Delete(ctx context.Context, short string) error {
	fmt.Printf("Deleting short URL key: %s\n", short)

	err := r.client.Del(ctx, short).Err()
	if err != nil {
		return fmt.Errorf("failed to delete short URL: %v", err)
	}

	fmt.Println("Deleted successfully.")
	return nil
}

func (r *RedisStore) All(ctx context.Context) (map[string]string, error) {
	result := make(map[string]string)

	// Use Redis SCAN to iterate keys safely
	iter := r.client.Scan(ctx, 0, "*", 0).Iterator()
	for iter.Next(ctx) {
		key := iter.Val()
		val, err := r.client.Get(ctx, key).Result()
		if err != nil && err != redis.Nil {
			return nil, fmt.Errorf("failed to get value for key %s: %v", key, err)
		}
		result[key] = val
	}
	if err := iter.Err(); err != nil {
		return nil, fmt.Errorf("error scanning keys: %v", err)
	}
	return result, nil
}
