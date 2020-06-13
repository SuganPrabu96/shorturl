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
	address string
	conn    redis.Conn
}

// NewClient creates a new redis client
func NewClient(host string, port int) (*Client, error) {
	ctx := context.Background()
	// TODO: Setup a connection pool
	address := host + ":" + strconv.Itoa(port)
	options := []redis.DialOption{redis.DialReadTimeout(ReadTimeout), redis.DialWriteTimeout(WriteTimeout)}
	c, err := redis.DialContext(ctx, "tcp", address, options...)
	if err != nil {
		return nil, err
	}
	defer c.Close()
	client := &Client{address, c}
	return client, nil
}

// GetValue gets value by key
func (c *Client) GetValue(key string) (string, error) {
	cmd := "GET " + key
	reply, err := c.conn.Do(cmd)
	if err != nil {
		return "", failure.Wrap(err)
	}
	if reply == nil {
		return "", nil
	}
	return fmt.Sprintf("%v", reply), nil
}

// SetValueIfNotExists sets value if it does not exist
func (c *Client) SetValueIfNotExists(key, value string) error {
	cmd := "SETNX" + key + " " + value
	reply, err := c.conn.Do(cmd)
	if err != nil {
		return failure.Wrap(err)
	}
	status := fmt.Sprintf("%v", reply)
	if status != "OK" {
		return errors.New("Key already exists")
	}
	return nil
}

// Close closes connection to redis server
func (c *Client) Close() error {
	err := c.conn.Close()
	if err != nil {
		return failure.Wrap(err)
	}
	return nil
}
