package main

import (
	"fmt"
	"os"

	"github.com/SuganPrabu96/shorturl/redis"
	"github.com/SuganPrabu96/shorturl/server"
	"github.com/kelseyhightower/envconfig"
	"github.com/morikuni/failure"
)

// Config is the application configuration
type Config struct {
	RedisHost string `envconfig:"REDIS_HOST" default:"localhost"`
	RedisPort int    `envconfig:"REDIS_PORT" default:"6379"`
}

func readConfig() (*Config, error) {
	var config Config
	err := envconfig.Process("", &config)
	if err != nil {
		return nil, failure.Wrap(err)
	}

	return &config, nil
}

func main() {
	c, err := readConfig()
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to read config: %s\n", err)
		os.Exit(1)
	}
	redisClient, err := redis.NewClient(c.RedisHost, c.RedisPort)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Failed to connect to redis server: %s\n", err)
		os.Exit(1)
	}

	s := server.NewServer(redisClient)
	s.Serve()
}
