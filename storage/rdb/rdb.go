/*
pacakge redis 包裝了 http server 中常用的方法，你可以用注入 redis rdb 的方式來使用這些方法，
使用注入而不包成 struct 的原因是為了保留操作上的彈性，以免增加少數特例條件時還要來 core 新增方法。
*/
package rdb

import (
	"fmt"
	"time"

	"github.com/win30221/core/http/catch"
	"github.com/win30221/core/http/ctx"
	"github.com/win30221/core/syserrno"

	jsoniter "github.com/json-iterator/go"
	"github.com/redis/go-redis/v9"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func Get(ctx *ctx.Context, rdb *redis.Client, key string, result any) (err error) {
	b, err := rdb.Get(ctx.Context, key).Bytes()
	if err == redis.Nil {
		return
	}
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "get data error", fmt.Sprintf("execute GET command error. err: %s", err.Error()), 3)
		return
	}

	err = json.Unmarshal(b, result)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "get data error", fmt.Sprintf("unmarshal data error. err: %s, got: %s", err.Error(), string(b)), 3)
		return
	}

	return
}

func Set(ctx *ctx.Context, rdb *redis.Client, key string, data any) (err error) {
	b, err := json.Marshal(data)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "set data error", fmt.Sprintf("marshal data error. err: %s, data: %+v", err.Error(), data), 3)
		return
	}

	_, err = rdb.Set(ctx.Context, key, b, 0).Result()
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "set data error", fmt.Sprintf("execute SET command error. err: %s", err.Error()), 3)
		return
	}

	return
}

func SetEX(ctx *ctx.Context, rdb *redis.Client, key string, ttl time.Duration, data any) (err error) {
	b, err := json.Marshal(data)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "setex data error", fmt.Sprintf("marshal data error. err: %s, data: %+v", err.Error(), data), 3)
		return
	}

	_, err = rdb.Set(ctx.Context, key, b, ttl).Result()
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "setex data error", fmt.Sprintf("execute SETEX command error. err: %s", err.Error()), 3)
		return
	}

	return
}

func SetNX(ctx *ctx.Context, rdb *redis.Client, key string, ttl time.Duration, data any) (err error) {
	b, err := json.Marshal(data)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "setnx data error", fmt.Sprintf("marshal data error. err: %s, data: %+v", err.Error(), data), 3)
		return
	}

	ok, err := rdb.SetNX(ctx.Context, key, b, ttl).Result()
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "setnx data error", fmt.Sprintf("execute SETNX command error. err: %s", err.Error()), 3)
		return
	}

	if !ok {
		err = catch.NewWitStack(syserrno.Redis, "setnx data error", "execute SETNX command error.", 3)
		return
	}

	return
}

func Del(ctx *ctx.Context, rdb *redis.Client, key string) (err error) {
	_, err = rdb.Del(ctx.Context, key).Result()
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "del error", fmt.Sprintf("execute DEL command error. err: %s", err.Error()), 3)
		return
	}

	return
}

func Exists(ctx *ctx.Context, rdb *redis.Client, key string) (isExist bool, err error) {
	res, err := rdb.Exists(ctx.Context, key).Result()
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "check exist error", fmt.Sprintf("execcute EXISTS command error. err: %s", err.Error()), 3)
		return
	}

	isExist = (res == 1)

	return
}

func FuzzyDel(ctx *ctx.Context, rdb *redis.Client, key string) (err error) {
	var cursor uint64
	var target []string
	var keys []string
	for {
		keys, cursor, err = rdb.Scan(ctx.Context, cursor, "*"+key+"*", 0).Result()
		if err != nil {
			err = catch.NewWitStack(syserrno.Redis, "fuzzy del error", fmt.Sprintf("execute SCAN command error. err: %s", err.Error()), 3)
			return
		}

		target = append(target, keys...)

		if cursor == 0 {
			break
		}
	}

	if len(target) == 0 {
		return
	}

	_, err = rdb.Unlink(ctx.Context, target...).Result()
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "fuzzy del error", fmt.Sprintf("execute UNLINK command error. err: %s", err.Error()), 3)
		return
	}

	return
}

func Lock(ctx *ctx.Context, rdb *redis.Client, key string, ttl time.Duration) (err error) {
	err = SetNX(ctx, rdb, key, ttl, "")
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "Lock error", fmt.Sprintf("Lock error. err: %s", err.Error()), 3)
		return
	}

	return
}

func Unlock(ctx *ctx.Context, rdb *redis.Client, key string) (err error) {
	err = Del(ctx, rdb, key)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "Unlock error", fmt.Sprintf("Unlock error. err: %s", err.Error()), 3)
		return
	}

	return
}
