package basic

import (
	"fmt"
	"log"

	"github.com/win30221/core/config"
)

func Init(serverName string) {
	ServerName = serverName
	// 收命令列參數
	loadEnv()
	// 初始化 Consul
	config.Load(ConsulIP)
	// Load consul env
	// 載入內部系統 Private Token。這個參數在 http middleware 的 valid_token 會使用到
	SysToken, _ = config.GetString("/system/systoken", true)
	Site, _ = config.GetString("/system/site", true)
	LogMode, _ = config.GetString(fmt.Sprintf("/service/%s/log_mode", ServerName), false)
	if LogMode == "" {
		LogMode, _ = config.GetString("/system/log_mode", true)
	}

	RequestLatencyThrottle, _ = config.GetMillisecond("/system/request_latency_throttle", true)

	// 設定 Log
	setLog()
	// 檢查必要 loading 的參數
	check()
	log.Println("Environment: " + Site)
}

// 檢查必要相依參數
func check() {
	if ServerName == "" {
		log.Fatal("Server name is empty")
	}
}
