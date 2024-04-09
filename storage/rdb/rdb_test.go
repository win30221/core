package rdb

import (
	"testing"
	"time"

	"github.com/redis/go-redis/v9"
	"github.com/win30221/core/config"
	"github.com/win30221/core/http/ctx"
	"github.com/win30221/core/storage"
	"github.com/win30221/core/utils"
)

func init() {
	config.Load("127.0.0.1")
}

func Test_Bulk(t *testing.T) {
	rdb := storage.GetRedis("/storage/redis/auction-master", "workers")

	c := ctx.NewEmpty()

	err := SetEX(c, rdb, "test", time.Second, "test")
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	ans := ""
	err = Get(c, rdb, "test", &ans)
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if ans != "test" {
		t.Errorf("result : %+v, ", ans)
		return
	}

	isExist, err := Exists(c, rdb, "test")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if !isExist {
		t.Errorf("result : %+v, ", isExist)
		return
	}

	time.Sleep(2 * time.Second)

	err = Get(c, rdb, "test", &ans)
	if err != nil && err != redis.Nil {
		t.Errorf("%v", err)
		return
	}

	isExist, err = Exists(c, rdb, "test")
	if err != nil {
		t.Errorf("%v", err)
		return
	}
	if isExist {
		t.Errorf("result : %+v, ", isExist)
		return
	}

}

func Test_FuzzyDel(t *testing.T) {
	redis := storage.GetRedis("/storage/redis/auction-master", "workers")

	testCase := []string{"TEST", "TEST:01", "test:02", "TE00ST:03", "00TEST:04", "TEST00:05"}

	c := ctx.NewEmpty()
	for _, v := range testCase {
		Set(c, redis, v, v)
	}

	err := FuzzyDel(c, redis, "TEST")
	if err != nil {
		t.Errorf("TEST : %v", err)
		return
	}

	//預期結果 test:02，TE00ST:03 ，共2筆
	var result []string
	var data string
	for _, v := range testCase {
		err = Get(c, redis, v, &data)
		if err != nil {
			continue
		}
		result = append(result, data)
	}

	if len(result) != 2 {
		t.Errorf("result : %+v, input count: %d  , ouput count:  %d", result, len(testCase), len(result))
	}
}

func Test_Lock(t *testing.T) {
	rdb := storage.GetRedis("/storage/redis/auction-master", "workers")

	c := ctx.NewEmpty()
	key := utils.GenerateRequestID()

	err := Lock(c, rdb, key, time.Second)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	err = Lock(c, rdb, key, time.Second)
	if err == nil {
		t.Errorf("%v", err)
		return
	}

	time.Sleep(2 * time.Second)
	err = Lock(c, rdb, key, time.Second)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	err = Unlock(c, rdb, key)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

	err = Lock(c, rdb, key, time.Second)
	if err != nil {
		t.Errorf("%v", err)
		return
	}

}
