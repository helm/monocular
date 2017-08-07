package config

import (
	log "github.com/Sirupsen/logrus"
	"github.com/albrow/zoom"
)

// Pool is a pool of Zoom connections used by other packages
var Pool *zoom.Pool

type redisConfig struct {
	Host string
}

// NewRedisPool initializes the pool of Zoom connections
func NewRedisPool() *zoom.Pool {
	config := getRedisConf()
	Pool = zoom.NewPool(config.Host)
	return Pool
}

func getRedisConf() redisConfig {
	config, err := GetConfig()
	if err != nil {
		log.Fatalf("unable to read config")
	}
	redis := config.Redis
	// Set default Redis host
	if redis.Host == "" {
		redis.Host = "localhost:6379"
	}
	return redis
}
