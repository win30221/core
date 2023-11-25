// See https://blog.twitter.com/engineering/en_us/a/2010/announcing-snowflake for more information.
package utils

import (
	"fmt"
	"math/big"
	"os"
	"sync"
	"time"
)

var (
	seq uint32
	mux sync.Mutex
)

// snowflaker
func GenerateSnowflakeID() string {
	/*
		Unix timestamp_PID_seq		1545112028_12345_12345
									15451120281234512345
		Maximum uint64				18446744073709551615
	*/
	mux.Lock()
	s := fmt.Sprintf("%d%d%d", time.Now().Unix(), os.Getpid(), seq)
	seq++
	mux.Unlock()

	i := new(big.Int)
	fmt.Sscan(s, i)

	return i.Text(62)
}
