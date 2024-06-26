package storage

import (
	"context"
	"log"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/spf13/cast"
	"github.com/win30221/core/config"
)

func GetRedis(path, dbName string) (rdb *redis.Client) {
	var err error

	defer func() {
		if err != nil {
			log.Fatalf("get redis error: %s \n - path %s - dbName %s", err, path, dbName)
		}
	}()

	// 先抓指定連線資訊
	m, _ := config.GetStringMap(path+"/"+dbName, false)
	host := m["host"]
	password := m["password"]
	maxIdle := cast.ToInt(m["max_idle"])
	maxActive := cast.ToInt(m["max_active"])
	idleTimeout, _ := time.ParseDuration(m["idle_timeout"])

	// 沒有設定再抓預設
	if host == "" {
		host, _ = config.GetString(path+"/host", true)
	}
	if password == "" {
		password, _ = config.GetString(path+"/password", true)
	}
	if maxIdle == 0 {
		maxIdle, _ = config.GetInt(path+"/max_idle", true)
	}
	if maxActive == 0 {
		maxActive, _ = config.GetInt(path+"/max_active", true)
	}
	if idleTimeout == 0 {
		idleTimeout, _ = config.GetDuration(path+"/idle_timeout", true)
	}

	db, _ := config.GetInt(path+"/dbname."+dbName, false)

	rdb = redis.NewClient(&redis.Options{
		Addr:         host,
		Password:     password,
		DB:           db,
		MaxIdleConns: maxIdle,
	})

	// rdb = &redigo.Pool{
	// 	MaxIdle:     maxIdle,
	// 	MaxActive:   maxActive,
	// 	Wait:        true,
	// 	IdleTimeout: idleTimeout,
	// 	Dial: func() (conn redigo.Conn, err error) {
	// 		conn, err = redigo.Dial("tcp", host)
	// 		if err != nil {
	// 			return
	// 		}

	// 		if password != "" {
	// 			_, err = conn.Do("AUTH", password)
	// 			if err != nil {
	// 				return
	// 			}
	// 		}

	// 		_, err = conn.Do("SELECT", db)
	// 		if err != nil {
	// 			conn.Close()
	// 			return
	// 		}

	// 		return
	// 	},
	// 	TestOnBorrow: func(c redigo.Conn, t time.Time) (err error) {
	// 		_, err = c.Do("PING")
	// 		return
	// 	},
	// }

	_, err = rdb.Ping(context.Background()).Result()

	log.Printf("Redis connected to `%+v`, selected db to `%+v` success", host, db)

	return
}
