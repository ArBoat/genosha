package db

import (
	"genosha/utils/myLogger"
	"go.uber.org/zap"

	"github.com/gomodule/redigo/redis"
)

var Redis redis.Conn

func init() {
	//Redis = createRedisHandler()
}

func createRedisHandler() redis.Conn {
	c, err := redis.Dial("tcp", "127.0.0.1:6379", redis.DialDatabase(9))
	if err != nil {
		myLogger.Log.Error("Connect to redis error", zap.Error(err))
		return nil
	}
	return c
}
