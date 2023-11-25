package basic

import (
	"flag"
	"log"
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
	Location   string
	TimeZone   *time.Location
	ServerName string
	// LogMode 允許參數 debug, info, warn, error, dpanic, panic, fatal
	LogMode string

	// Alert
	// RequestLatencyThrottle 如果請求時長大於 RequestLatencyThrottle 時需要印出 warn
	RequestLatencyThrottle int
)

func loadEnv() {
	flag.StringVar(&ConsulIP, "c", "127.0.0.1", "Consul IP")
	flag.StringVar(&Host, "h", "0.0.0.0", "Server Host")
	flag.StringVar(&Port, "p", "1324", "Server Port")
	flag.StringVar(&Location, "l", "Asia/Taipei", "Time zone")
	flag.Parse()

	timeZone, err := time.LoadLocation(Location)
	if err != nil {
		log.Fatalln(err)
	}

	TimeZone = timeZone
	time.Local = timeZone
}
