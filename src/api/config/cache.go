package config

import (
	"sync"

	log "github.com/Sirupsen/logrus"
	"github.com/albrow/zoom"
)

const defaultHost = "localhost:6379"

// Pool is a pool of Zoom connections used by other packages
var pool *zoom.Pool

var once sync.Once

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

// CloseRedisPool closes a pool of Zoom connections
func CloseRedisPool() {
	if pool == nil {
		return
	}
	if err := pool.Close(); err != nil {
		log.Fatalf("unable to close pool")
	}
	pool = nil
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
