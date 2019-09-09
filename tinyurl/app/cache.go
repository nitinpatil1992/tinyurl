package app

import (
	"strconv"
	"time"

	"github.com/go-redis/redis"
	log "github.com/sirupsen/logrus"
)

type Result struct {
	ShortString string
	CreatedAt   time.Time
}

func GetRedisClient(redisHost string) *redis.Client {
	client := redis.NewClient(&redis.Options{
		Addr: redisHost,
	})
	return client
}

func getCachedURL(longURL string) (bool, Result) {

	var result Result
	cacheExists, _ := redisClient.Exists(longURL).Result()

	if cacheExists == 0 {
		return false, result
	}

	log.Info("Cache hit occurred")

	shortString, _ := redisClient.HGet(longURL, "short_string").Result()
	createdAtStringTime, _ := redisClient.HGet(longURL, "created_at").Result()
	createdAtUnixTime, err := strconv.ParseInt(createdAtStringTime, 10, 64)

	if err != nil {
		log.Warn("failed to convert string time to unix time")
		return false, result
	}

	createdAt := time.Unix(createdAtUnixTime, 0)

	result = Result{
		ShortString: shortString,
		CreatedAt:   createdAt,
	}
	return true, result
}

func setCacheURL(longURL, shortString string, currentTime time.Time) {
	_ = redisClient.HSet(longURL, "short_string", shortString)
	_ = redisClient.HSet(longURL, "created_at", currentTime.Unix())
}
