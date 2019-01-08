package main

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

var _ = redis.Dial

var client = &redis.Pool{
	MaxIdle:     10,
	MaxActive:   60,
	IdleTimeout: 180 * time.Second,
	Dial: func() (redis.Conn, error) {
		c, err := redis.Dial("tcp", "127.0.0.1:6379")
		if err != nil {
			return nil, err
		}

		c.Do("SELECT", 0)
		return c, nil
	},
}

func main() {
	conn := client.Get()
	defer conn.Close()

	n, err := conn.Do("SET", "testkey", "testvalue")
	if err != nil {
		fmt.Println("Do.SET: ", err)
		return
	}
	fmt.Println(n)

	fmt.Println("main")
}
