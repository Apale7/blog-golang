package common

import (
	"github.com/garyburd/redigo/redis"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
	"time"
)

var (
	redisPool *redis.Pool
)

func init() {
	redisPool = redisPoolInit("localhost:6379", "")
}

func redisPoolInit(server string, password string) *redis.Pool {
	return &redis.Pool{
		MaxIdle:     50,
		MaxActive:   30,
		IdleTimeout: 300 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", server)
			if err != nil {
				return nil, err
			}
			if password != "" {
				if _, err := c.Do("AUTH", password); err != nil {
					_ = c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}

func RedisExec(cmd string, key interface{}, args ...interface{}) (interface{}, error) {
	con := redisPool.Get()
	if err := con.Err(); err != nil {
		return nil, err
	}
	defer con.Close()
	parmas := make([]interface{}, 0)
	parmas = append(parmas, key)

	if len(args) > 0 {
		for _, v := range args {
			parmas = append(parmas, v)
		}
	}
	return con.Do(cmd, parmas...)
}

func Exist(key string) bool  {
	c := redisPool.Get()
	isKeyExist, err := redis.Bool(c.Do("EXISTS", key))
	if err != nil {
		log.Error(errors.WithStack(err))
		panic(errors.WithStack(err))
	}
	return isKeyExist
}