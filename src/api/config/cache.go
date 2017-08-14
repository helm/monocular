package config

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/albrow/zoom"
)

const defaultHost = "localhost:6379"

var (
	pool *zoom.Pool
	once sync.Once
)

type redisConfig struct {
	Host string
}

// GetRedisPool returns a pool of Zoom connections
func GetRedisPool() *zoom.Pool {
	once.Do(func() {
		pool = newRedisPool()
	})
	return pool
}

func newRedisPool() *zoom.Pool {
	config := getRedisConf()
	return zoom.NewPool(config.Host)
}

func getRedisConf() redisConfig {
	config, err := GetConfig()
	if err != nil {
		log.Fatalf("unable to read config")
	}
	redis := config.Redis
	// Set default Redis host
	if redis.Host == "" {
		redis.Host = defaultHost
	}
	return redis
}
