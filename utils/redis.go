package utils

import (
	"fmt"
	"time"

	"github.com/gomodule/redigo/redis"
)

// interface
type RedisPool interface {
	InitRedis(v ValidateInfo)
	GetCon() redis.Conn
}

// redis pool
type RedisPoolInstance struct {
	_g_pool *redis.Pool
}

// 初始化
func (r *RedisPoolInstance) InitRedis(v ValidateInfo) {
	if r._g_pool != nil {
		return
	}
	if v.PassWord != "" {
		r._g_pool = &redis.Pool{
			MaxIdle:     256,
			MaxActive:   0,
			IdleTimeout: time.Duration(120),
			Dial: func() (redis.Conn, error) {
				return redis.Dial(
					"tcp",
					fmt.Sprintf("%s:%d", v.Host, v.Port),
					redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
					redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
					redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
					redis.DialDatabase(0),
					redis.DialPassword(v.PassWord),
				)
			},
		}
	} else {
		r._g_pool = &redis.Pool{
			MaxIdle:     256,
			MaxActive:   0,
			IdleTimeout: time.Duration(120),
			Dial: func() (redis.Conn, error) {
				return redis.Dial(
					"tcp",
					fmt.Sprintf("%s:%d", v.Host, v.Port),
					redis.DialReadTimeout(time.Duration(1000)*time.Millisecond),
					redis.DialWriteTimeout(time.Duration(1000)*time.Millisecond),
					redis.DialConnectTimeout(time.Duration(1000)*time.Millisecond),
					redis.DialDatabase(0),
				)
			},
		}

	}
}

// 获取链接
func (r *RedisPoolInstance) GetCon() redis.Conn {
	return r._g_pool.Get()
}
