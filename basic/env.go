package basic

import (
	"flag"
	"time"
)

var (
	Port       string
	Host       string
	SysToken   string
	ConsulIP   string
	Debug      bool
	Version    string
	Commit     string
	BuildTime  string
	Site       string
	TimeZone   = time.FixedZone("UTC+9", 9*60*60)
	ServerName string
	// LogMode 允許參數 debug, info, warn, error, dpanic, panic, fatal
	LogMode string

	// Alert
	// RequestLatencyThrottle 如果請求時長大於 RequestLatencyThrottle 時需要印出 warn
	RequestLatencyThrottle int
)

func loadEnv() {
	flag.StringVar(&ConsulIP, "c", "127.0.0.1", "Consul IP")
	flag.StringVar(&Port, "p", "1324", "Server Port")
	flag.StringVar(&Host, "b", "0.0.0.0", "Server Host")
	flag.Parse()

	time.Local = TimeZone
}
