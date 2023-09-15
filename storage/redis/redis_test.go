package redis

import (
	"context"
	"testing"

	"github.com/win30221/core/config"
	"github.com/win30221/core/storage"
)

func init() {
	config.Load("127.0.0.1")
}

func Test_FuzzyDel(t *testing.T) {
	redisPool := storage.GetRedis("/storage/redis/promote-ms", "user")

	testCase := []string{"TEST", "TEST:01", "test:02", "TE00ST:03", "00TEST:04", "TEST00:05"}

	ctx := context.TODO()
	for _, v := range testCase {
		SET(redisPool, ctx, v, v)
	}

	err := FuzzyDel(redisPool, ctx, "TEST")
	if err != nil {
		t.Errorf("TEST : %v", err)
		return
	}

	//預期結果 test:02，TE00ST:03 ，共2筆
	var result []string
	var data string
	for _, v := range testCase {
		err = GET(redisPool, ctx, v, &data)
		if err != nil {
			continue
		}
		result = append(result, data)
	}

	if len(testCase) != len(result) {
		t.Errorf("result : %+v, input count: %d  , ouput count:  %d", result, len(testCase), len(result))
	}

	return
}
