package redis

import (
    "fmt"
    "github.com/gomodule/redigo/redis"
    jsoniter "github.com/json-iterator/go"
    "magpie-gateway/configuration"
    "time"
)

/*
 Be careful to use, this library will not prevent sql inject.
*/

var (
    redisPool *redis.Pool
)


func newPool(uri string) *redis.Pool {
    return &redis.Pool{
        Dial: func() (redis.Conn, error) {
            return redis.DialURL(uri)
        },
        MaxIdle:     10,
        IdleTimeout: 180 * time.Second,
        TestOnBorrow: func(c redis.Conn, t time.Time) error {
            if time.Since(t) < time.Minute {
                return nil
            }
            _, err := c.Do("PING")
            return err
        },
    }
}

func GetConn() redis.Conn {
    return redisPool.Get()
}

func SetString(key, value string, options string) error {
    conn := GetConn()
    _, err := conn.Do(fmt.Sprintf("SET \"%s\" \"%s\" %s", key, value, options))
    return err
}

func GetString(key string) (string, error) {
    conn := GetConn()
    res, err := conn.Do(fmt.Sprintf("GET \"%s\"", key))
    if err != nil {
        return "", err
    }
    r, ok := res.(string)
    if !ok {
        return "", redis.Error("Could not parse key to string")
    }
    return r, nil
}

func SetList(key, options string, value ...interface{}) error {
    conn := GetConn()

    for i := range value {
        data, err := jsoniter.Marshal(&value[i])
        if err != nil {
            return nil
        }
        if _, err = conn.Do(fmt.Sprintf("LPUSH \"%s\" \"%s\" %s", key, string(data), options)); err != nil {
            return err
        }
    }
    return nil
}

func init() {
    redisPool = newPool(configuration.GlobalConfiguration.RedisURI)
}
