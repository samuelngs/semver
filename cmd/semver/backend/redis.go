package backend

import (
	"os"

	"gopkg.in/redis.v3"
)

var client *redis.Client

func init() {
	opts := &redis.Options{
		Addr:       "localhost:6379",
		DB:         0,
		MaxRetries: 5,
	}
	if os.Getenv("DOCKER") != "" {
		opts.Addr = "redis:6379"
	}
	client = redis.NewClient(opts)
}

// Exists checks if path or dir exists in database
func Exists(path string) (bool, error) {
	_, err := client.Get(path).Result()
	if err == redis.Nil {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

// Get data from database and return result
func Get(path string) (string, error) {
	val, err := client.Get(path).Result()
	if err != nil {
		return "", err
	}
	return val, nil
}

// GetAll data from database and return result
func GetAll(path ...string) ([]interface{}, error) {
	arr, err := client.MGet(path...).Result()
	if err != nil {
		return nil, err
	}
	return arr, nil
}

// Set data to database and return result
func Set(val string, dirs ...string) error {
	for _, dir := range dirs {
		if err := client.Set(dir, val, 0).Err(); err != nil {
			return err
		}
	}
	return nil
}

// Delete data from database
func Delete(dirs ...string) error {
	for _, dir := range dirs {
		if err := client.Del(dir).Err(); err != nil {
			return err
		}
	}
	return nil
}

// List data from database and return result
func List(path string) ([]string, error) {
	return client.Keys(path).Result()
}
