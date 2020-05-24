package dao

import (
	"encoding/json"
	"genosha/db"
	"genosha/utils/myLogger"
	"github.com/gomodule/redigo/redis"
	"go.uber.org/zap"
)

var c = db.Redis

func init() {
	//testRedis()
}

func testRedis() {
	key := "profile"
	imap := map[string]string{"username": "666", "phonenumber": "888"}
	value, _ := json.Marshal(imap)

	n, err := c.Do("SETNX", key, value)
	if err != nil {
		myLogger.Log.Error("redis error", zap.Error(err))
	}
	if n == int64(1) {
		myLogger.Log.Info("success")
	}

	var imapGet map[string]string

	valueGet, err := redis.Bytes(c.Do("GET", key))
	if err != nil {
		myLogger.Log.Error("redis error", zap.Error(err))
	}

	errShal := json.Unmarshal(valueGet, &imapGet)
	if errShal != nil {
		myLogger.Log.Error("redis error", zap.Error(err))
	}
	myLogger.Log.Info(imapGet["username"])
	myLogger.Log.Info(imapGet["phonenumber"])
}

func RedisKeyExist(key string) (bool, error) {
	exists, err := redis.Bool(c.Do("EXISTS", key))
	return exists, err
}
func RedisGetKey(key string) (string, error) {
	v, err := redis.String(c.Do("GET", key))
	return v, err
}
func RedisSetKeyWithExpireTime(key, value, secondSting string) error {
	_, err := c.Do("SET", key, value, "EX", secondSting)
	return err
}
