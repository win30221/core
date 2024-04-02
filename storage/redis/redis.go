/*
pacakge redis 包裝了 http server 中常用的方法，你可以用注入 redis pool 的方式來使用這些方法，
使用注入而不包成 struct 的原因是為了保留操作上的彈性，以免增加少數特例條件時還要來 core 新增方法。
*/
package redis

import (
	"context"
	"fmt"

	"github.com/win30221/core/http/catch"
	"github.com/win30221/core/syserrno"

	"github.com/gomodule/redigo/redis"
	redigo "github.com/gomodule/redigo/redis"
	jsoniter "github.com/json-iterator/go"
)

var json = jsoniter.ConfigCompatibleWithStandardLibrary

func GET(pool *redigo.Pool, ctx context.Context, key string, result interface{}) (err error) {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "get data error", fmt.Sprintf("get connection error. err: %s", err.Error()), 3)
		return
	}
	defer conn.Close()

	b, err := redigo.Bytes(conn.Do("GET", key))
	if err == redis.ErrNil {
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

func SET(pool *redigo.Pool, ctx context.Context, key string, data interface{}) (err error) {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "set data error", fmt.Sprintf("get connection error. err: %s", err.Error()), 3)
		return
	}
	defer conn.Close()

	b, err := json.Marshal(data)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "set data error", fmt.Sprintf("marshal data error. err: %s, data: %+v", err.Error(), data), 3)
		return
	}

	_, err = conn.Do("SET", key, string(b))
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "set data error", fmt.Sprintf("execute SET command error. err: %s", err.Error()), 3)
		return
	}

	return
}

func SETEX(pool *redigo.Pool, ctx context.Context, key string, ttl int, data interface{}) (err error) {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "set data error", fmt.Sprintf("get connection error. err: %s", err.Error()), 3)
		return
	}
	defer conn.Close()

	b, err := json.Marshal(data)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "set data error", fmt.Sprintf("marshal data error. err: %s, data: %+v", err.Error(), data), 3)
		return
	}

	_, err = conn.Do("SETEX", key, ttl, string(b))
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "set data error", fmt.Sprintf("execute SETEX command error. err: %s", err.Error()), 3)
		return
	}

	return
}

func DEL(pool *redigo.Pool, ctx context.Context, key string) (err error) {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "del error", fmt.Sprintf("get connection error. err: %s", err.Error()), 3)
		return
	}
	defer conn.Close()

	_, err = conn.Do("DEL", key)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "del error", fmt.Sprintf("execute DEL command error. err: %s", err.Error()), 3)
		return
	}

	return
}

func FuzzyDel(pool *redigo.Pool, ctx context.Context, key string) (err error) {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "fuzzy del error", fmt.Sprintf("get connection error. err: %s", err.Error()), 3)
		return
	}
	defer conn.Close()

	iter := 0
	var target []interface{}
	var arr []interface{}
	for {

		arr, err = redis.Values(conn.Do("SCAN", iter, "MATCH", "*"+key+"*"))
		if err != nil {
			err = catch.NewWitStack(syserrno.Redis, "fuzzy del error", fmt.Sprintf("execute SCAN command error. err: %s", err.Error()), 3)
			return
		}

		iter, err = redis.Int(arr[0], nil)
		if err != nil {
			err = catch.NewWitStack(syserrno.Redis, "fuzzy del error", fmt.Sprintf("redis INT function error. err: %s", err.Error()), 3)
			return
		}

		keys, err := redis.Values(arr[1], nil)
		if err != nil {
			err = catch.NewWitStack(syserrno.Redis, "fuzzy del error", fmt.Sprintf("redis Strings function error. err: %s", err.Error()), 3)
			return err
		}
		target = append(target, keys...)

		if iter == 0 {
			break
		}
	}
	if len(target) == 0 {
		return
	}

	_, err = conn.Do("UNLINK", target...)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "fuzzy del error", fmt.Sprintf("execute UNLINK command error. err: %s", err.Error()), 3)
		return
	}

	return
}

func EXISTS(pool *redigo.Pool, ctx context.Context, key string) (isExist bool, err error) {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "check exist error", fmt.Sprintf("get connection error. err: %s", err.Error()), 3)
		return
	}
	defer conn.Close()

	isExist, err = redis.Bool(conn.Do("EXISTS", key))
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "check exist error", fmt.Sprintf("execcute EXISTS command error. err: %s", err.Error()), 3)
		return
	}

	return
}

func SADD(pool *redigo.Pool, ctx context.Context, key string, data ...interface{}) (err error) {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "SADD data error", fmt.Sprintf("get connection error. err: %s", err.Error()), 3)
		return
	}
	defer conn.Close()

	args := []interface{}{key}
	args = append(args, data...)

	_, err = conn.Do("SADD", args...)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "SADD data error", fmt.Sprintf("execute SADD command error. err: %s", err.Error()), 3)
		return
	}

	return
}

func SMEMBERS(pool *redigo.Pool, ctx context.Context, key string, result *[]string) (err error) {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "SMEMBERS data error", fmt.Sprintf("get connection error. err: %s", err.Error()), 3)
		return
	}
	defer conn.Close()

	list, err := redis.Strings(conn.Do("SMEMBERS", key))
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "SMEMBERS data error", fmt.Sprintf("execute SMEMBERS command error. err: %s", err.Error()), 3)
		return
	}

	*result = append(*result, list...)

	return
}

func Lock(pool *redigo.Pool, ctx context.Context, key string, ttl int) (err error) {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "Lock error", fmt.Sprintf("Lock error. err: %s", err.Error()), 3)
		return
	}
	defer conn.Close()

	reply, err := redis.String(conn.Do("SET", key, 1, "NX", "PX", ttl))
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "Lock error", fmt.Sprintf("Lock error. err: %s", err.Error()), 3)
		return
	}

	if reply != "OK" {
		err = catch.NewWitStack(syserrno.Redis, "Lock error", "Lock error", 3)
		return
	}

	return
}

func Unlock(pool *redigo.Pool, ctx context.Context, key string) (err error) {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "Unlock error", fmt.Sprintf("get connection error. err: %s", err.Error()), 3)
		return
	}
	defer conn.Close()

	reply, err := redis.Int(conn.Do("DEL", key))
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "Unlock error", fmt.Sprintf("Unlock error. err: %s", err.Error()), 3)
		return
	}

	if reply != 1 {
		err = catch.NewWitStack(syserrno.Redis, "Unlock error", "Unlock error", 3)
		return
	}

	return
}
