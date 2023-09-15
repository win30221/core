package redis

import (
	"context"
	"fmt"
	"time"

	redigo "github.com/gomodule/redigo/redis"
	"github.com/win30221/core/http/catch"
	"github.com/win30221/core/syserrno"
	"go.uber.org/zap"
)

func PUBLISH(pool *redigo.Pool, ctx context.Context, channel string, data interface{}) (err error) {
	conn, err := pool.GetContext(ctx)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "Publish error", fmt.Sprintf("get connection error. err: %s", err.Error()), 3)
		return
	}
	defer conn.Close()

	_, err = conn.Do("PUBLISH", channel, data)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "Publish error", fmt.Sprintf("publish error. err: %s", err.Error()), 3)
		return
	}

	return
}

type onMessage func(pubSubConn redigo.PubSubConn)

// Subscribe
// f: function 內部請使用 for { switch... } 來實作才能持續接收訊息
func Subscribe(pool *redigo.Pool, onMessage onMessage, channel ...interface{}) {
	if len(channel) == 0 {
		zap.L().Warn("Can not subscribe redis without channel")
		return
	}

	zap.L().Info("Subscribe to redis channel.", zap.Any("channel", channel))

	// subscribe redis channel 及處理重新連線
	go func() {
		reconnectTimer := time.NewTimer(5 * time.Second)
		init := false

		for {
			if !init {
				subscribe(pool, onMessage, channel...)
				init = true
				continue
			}

			reconnectTimer.Reset(5 * time.Second)
			<-reconnectTimer.C

			err := subscribe(pool, onMessage, channel...)
			if err != nil {
				zap.L().Error("subscribe error: " + err.Error())
			}
		}
	}()
}

func subscribe(pool *redigo.Pool, onMessage onMessage, channel ...interface{}) (err error) {
	conn, err := pool.GetContext(context.Background())
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "Subscribe error", fmt.Sprintf("get connection error. err: %s", err.Error()), 3)
		return
	}
	defer conn.Close()

	psc := redigo.PubSubConn{Conn: conn}

	err = psc.Subscribe(channel...)
	if err != nil {
		err = catch.NewWitStack(syserrno.Redis, "Subscribe error", fmt.Sprintf("subscribe error. err: %s", err.Error()), 3)
		return
	}

	onMessage(psc)

	zap.L().Warn("Connection of redis channel has been closed", zap.Any("channel", channel))

	return
}
