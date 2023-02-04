package db

import (
	"github.com/BlsaCn/short-url/conf"
	"github.com/go-redis/redis"
)

func NewRedisCli() *redis.Client {
	r := &conf.RedisCfg{}
	r.Redis()
	c := redis.NewClient(&redis.Options{
		Addr:     r.Host,
		Password: r.Pwd,
		DB:       r.DB,
	})
	if _, err := c.Ping().Result(); err != nil {
		panic(err)
	}
	return c
}
