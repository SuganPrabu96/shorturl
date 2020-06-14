package redis

import (
	"context"
	"errors"
	"fmt"
	"strconv"
	"time"

	"github.com/gomodule/redigo/redis"
	"github.com/morikuni/failure"
)

const (
	// ReadTimeout is timeout to read from redis server
	ReadTimeout = time.Second
	// WriteTimeout is timeout to write to redis server
	WriteTimeout = time.Second
)

// Client is the client to connect to the redis server
type Client struct {
	Address string
	Pool    *redis.Pool
	Conn    redis.Conn
}

// NewClient creates a new redis client
func NewClient(host string, port int) (*Client, error) {
	ctx := context.Background()
	// TODO: Setup a connection pool
	address := host + ":" + strconv.Itoa(port)
	options := []redis.DialOption{redis.DialReadTimeout(ReadTimeout), redis.DialWriteTimeout(WriteTimeout)}

	pool := redis.Pool{
		// Maximum number of idle connections in the pool.
		MaxIdle: 80,
		// max number of connections
		MaxActive: 12000,
		// Dial is an application supplied function for creating and
		// configuring a connection.
		Dial: func() (redis.Conn, error) {
			c, err := redis.DialContext(ctx, "tcp", address, options...)
			if err != nil {
				return nil, err
			}
			return c, err
		},
	}

	client := &Client{address, &pool, nil}
	return client, nil
}

// GetValue gets value by key. Returns empty string if key was not found
func (c *Client) GetValue(key string) (string, error) {
	reply, err := redis.String(c.Conn.Do("GET", key))
	if err != nil {
		return "", failure.Wrap(err)
	}
	return fmt.Sprintf("%v", reply), nil
}

// SetValueIfNotExists sets value if it does not exist
func (c *Client) SetValueIfNotExists(key, value string) error {
	reply, err := c.Conn.Do("SETNX", key, value)
	if err != nil {
		return failure.Wrap(err)
	}
	status := reply.(int64)
	if status != 1 {
		return errors.New("Key already exists")
	}
	return nil
}
