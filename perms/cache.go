package perms

import (
	"fmt"
	"magpie-gateway/store/redis"
	"strconv"
)

var valPrefix = "perm::cache::"

func setCache(key string, value bool) error {
	return redis.SetString(fmt.Sprintf("%s%s", valPrefix, key), strconv.FormatBool(value), "")
}

func getCache(key string) (bool, error) {
	res, err := redis.GetString(fmt.Sprintf("%s%s", valPrefix, key))
	if err != nil {
		return false, err
	}
	r, err := strconv.ParseBool(res)
	if err != nil {
		return false, err
	}
	return r, nil
}