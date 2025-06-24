package store

import (
	"context"
	"os"
	"testing"

	"github.com/joho/godotenv"
)

func TestRedisStore_ConnectionSetGetClose(t *testing.T) {
	errLoad := godotenv.Load()
	if errLoad != nil {
		t.Fatal("Failed to load .env file:", errLoad)
	}

	if os.Getenv("REDIS_HOST") == "" {
		t.Fatal("REDIS_HOST not set")
	}

	ctx := context.Background() // top level context
	rdb := NewRedisStore()      // RedisStore* (which has a pointer to client)
	defer func() {
		if err := rdb.Close(); err != nil {
			t.Errorf("Failed to close Redis client: %v", err)
		}
	}()

	// Test SET
	key := "test:shortcode"
	value := "https://example.com"
	err := rdb.Set(ctx, key, value)
	if err != nil {
		t.Fatalf("Failed to set key: %v", err)
	}

	// Test GET
	got, err := rdb.Get(ctx, key)
	if err != nil {
		t.Fatalf("Failed to get key: %v", err)
	}
	if got != value {
		t.Errorf("Expected %s, got %s", value, got)
	}

	// Clean up
	if err := rdb.client.Del(ctx, key).Err(); err != nil {
		t.Logf("Warning: failed to delete test key: %v", err)
	}
}
