package lib

import (
	"errors"
	"log"
	"time"

	"github.com/go-redis/redis"
)

// RedisConn allows to create a connection with Redis storage
func RedisConn(dbSelected int) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",         // no password set
		DB:       dbSelected, // use default DB
	})
	err := client.Ping().Err()
	if err != nil {
		log.Panic(err)
	}
	return client
}

// RedisSetValue allows to set the value and expiration data of a key
func RedisSetValue(client *redis.Client, key string, val interface{}, expiration time.Duration) error {
	err := client.Set(key, val, expiration).Err()
	if err != nil {
		return err
	}
	return nil
}

// RedisGetValue allows to get the value of a key, return this value
func RedisGetValue(client *redis.Client, key string) (interface{}, error) {
	value, err := client.Get(key).Result()
	if err == redis.Nil {
		return nil, errors.New("Key does not exist")
	} else if err != nil {
		return nil, err
	}
	return value, nil
}

// RedisDelValue allows to remove a key
func RedisDelValue(client *redis.Client, key string) bool {
	ret := client.Del(key)
	if ret.Val() == 0 {
		return false
	}
	return true
}
