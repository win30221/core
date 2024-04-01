package config

import (
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/spf13/cast"
	"github.com/spf13/viper"

	_ "github.com/spf13/viper/remote"
)

var (
	ip                 = "127.0.0.1"
	ErrOnTypeIncorrect = errors.New("type incorrect")
)

func Ping() (err error) {
	method := "GET"
	url := "http://" + ip + ":8500/v1/status/peers"

	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return fmt.Errorf("resp code %+v", resp.StatusCode)
	}

	return
}

func Load(newIP string) {
	ip = newIP

	if err := Ping(); err != nil {
		log.Fatalf("Error on ping consul: %+v", err)
	}

	log.Println("Consul IP: " + ip)
}

func Get(key string, existOnErr bool, convert func(interface{}) error) (err error) {
	defer func() {
		if err != nil {
			if existOnErr {
				log.Fatalf("Error on load `%+v` from consul, Err: %v", key, err.Error())
			}

			log.Printf("Error on load `%+v` from consul, Err: %v", key, err.Error())
		}
	}()

	vObj := viper.New()

	// 取得 k, consul 設定檔中的 key 名稱
	// example:
	// input: "/storage/redis/config/account"
	// k = account
	sectionAry := strings.Split(key, "/")
	k := sectionAry[len(sectionAry)-1]

	// 取得 path
	// example:
	// input: "/storage/redis/config/account"
	// path = "/storage/redis/config"
	sectionAry = sectionAry[:len(sectionAry)-1]
	path := strings.Join(sectionAry, "/")

	vObj.AddRemoteProvider("consul", ip+":8500", path)
	vObj.SetConfigType("toml")
	err = vObj.ReadRemoteConfig()
	if err != nil {
		err = fmt.Errorf("%v (no section: %v)", err, key)
		return
	}

	res := vObj.Get(k)
	if res == nil {
		err = errors.New("Key `" + key + "` Not Found")
		return
	}

	err = convert(res)
	if err != nil {
		return
	}

	return
}

// GetString
// 假設 consul 路徑 "/storage/redis" 內有下列資料
// `
//
//	account = "hugo"
//	[dbname]
//	user = 2
//
// `
//
// 使用 GetStringMap("/storage/redis/dbname.user)
// return: "2"
// 使用 GetStringMap("/storage/redis/account)
// return: "hugo"
func GetString(key string, existOnErr bool) (result string, err error) {
	err = Get(key, existOnErr, func(res interface{}) (err error) {
		if _, ok := res.(string); !ok {
			err = ErrOnTypeIncorrect
			return
		}
		result = cast.ToString(res)
		return
	})
	return
}

func GetInt(key string, existOnErr bool) (result int, err error) {
	err = Get(key, existOnErr, func(res interface{}) (err error) {
		if _, ok := res.(int64); !ok {
			err = ErrOnTypeIncorrect
			return
		}
		result = cast.ToInt(res)
		return
	})
	return
}

func GetInt64(key string, existOnErr bool) (result int64, err error) {
	err = Get(key, existOnErr, func(res interface{}) (err error) {
		if _, ok := res.(int64); !ok {
			err = ErrOnTypeIncorrect
			return
		}
		result = cast.ToInt64(res)
		return
	})
	return
}

func GetFloat64(key string, existOnErr bool) (result float64, err error) {
	err = Get(key, existOnErr, func(res interface{}) (err error) {
		if _, ok := res.(float64); !ok {
			err = ErrOnTypeIncorrect
			return
		}
		result = cast.ToFloat64(res)
		return
	})
	return
}

func GetBool(key string, existOnErr bool) (result bool, err error) {
	err = Get(key, existOnErr, func(res interface{}) (err error) {
		if _, ok := res.(bool); !ok {
			err = ErrOnTypeIncorrect
			return
		}
		result = cast.ToBool(res)
		return
	})
	return
}

// GetStringSlice
// 假設 consul 路徑 "/storage/redis" 內有下列資料
// `
//
//	host = ["127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082"]
//	password = "b"
//
// `
//
// GetStringSlice("/storage/redis/host")
// return: []string{"127.0.0.1:8080", "127.0.0.1:8081", "127.0.0.1:8082"}
func GetStringSlice(key string, existOnErr bool) (result []string, err error) {
	err = Get(key, existOnErr, func(res interface{}) (err error) {
		if _, ok := res.([]interface{}); !ok {
			err = ErrOnTypeIncorrect
			return
		}
		result = cast.ToStringSlice(res)
		return
	})
	return
}

// GetDuration
// 假設 consul 路徑 "/storage/redis" 內有下列資料
// `
//
//	user_ttl = "60s"
//
// `
//
// GetDuration("/storage/redis/user_ttl")
// return: time.Duration("60s")
func GetDuration(key string, existOnErr bool) (result time.Duration, err error) {
	err = Get(key, existOnErr, func(res interface{}) (err error) {
		if _, ok := res.(string); !ok {
			err = ErrOnTypeIncorrect
			return
		}

		d := cast.ToString(res)
		result, err = time.ParseDuration(d)
		return
	})

	return
}

// GetSeconds
// 假設 consul 路徑 "/storage/redis" 內有下列資料
// `
//
//	user_ttl = "2m"
//
// `
//
// GetSeconds("/storage/redis/user_ttl")
// return: 120
func GetSeconds(key string, existOnErr bool) (result int, err error) {
	err = Get(key, existOnErr, func(res interface{}) (err error) {
		if _, ok := res.(string); !ok {
			err = ErrOnTypeIncorrect
			return
		}

		d := cast.ToString(res)
		parseDuration, err := time.ParseDuration(d)
		if err != nil {
			return
		}

		result = int(parseDuration.Seconds())

		return
	})

	return
}

// GetMillisecond
// 假設 consul 路徑 "/storage/redis" 內有下列資料
// `
//
//	user_ttl = "1s"
//
// `
//
// GetSeconds("/storage/redis/user_ttl")
// return: 1000
func GetMillisecond(key string, existOnErr bool) (result int, err error) {
	err = Get(key, existOnErr, func(res interface{}) (err error) {
		if _, ok := res.(string); !ok {
			err = ErrOnTypeIncorrect
			return
		}

		d := cast.ToString(res)
		parseDuration, err := time.ParseDuration(d)
		if err != nil {
			return
		}

		result = int(parseDuration.Milliseconds())

		return
	})

	return
}

// GetStringMap
// 假設 consul 路徑 "/storage/redis" 內有下列資料
// `
//
//	account = "hugo"
//	password = "b"
//	[dbname]
//	config = 1
//	user = 2
//
// `
//
// 如果要取得 dbname 的 map[string]string，可以使用 GetStringMap 方法，傳入參數為 GetStringMap("/storage/redis/dbname)
// return:
// `
//
//	map[string]string{
//		"config": 1,
//		"user": 2,
//	}
func GetStringMap(key string, existOnErr bool) (result map[string]string, err error) {
	err = Get(key, existOnErr, func(res interface{}) (err error) {
		result = cast.ToStringMapString(res)
		return
	})
	return
}

// GetFileStringMap
// ================= 備註 =================
// 目前架構中應該不會使用到這個方法，這是重構階段暫時留存用的，
// 通常使用 GetString/Get... 的方法及 GetStringMap 就能解決大部分情境。
// =======================================
// 假設 consul 路徑 "/storage/mongo" 內有下列資料
// `
//
//	account = "hugo"
//	password = "b"
//	[dbname]
//	config = 1
//	user = 2
//
// `
//
// 使用 GetStringMap 方法，傳入參數為 GetStringMap("/storage/mongo)
// return:
// `
//
//	map[string]string{
//		"account": "hugo",
//		"password": "b",
//	 "dbname.config": "1",
//	 "dbname.user": "2",
//	}
func GetFileStringMap(path string, existOnErr bool) (result map[string]string, err error) {
	defer func() {
		if err != nil {
			if existOnErr {
				os.Exit(1)
			}

			log.Printf("Error on load `%+v` from consul, Err: %v", path, err.Error())
		}
	}()

	vObj := viper.New()

	vObj.AddRemoteProvider("consul", ip+":8500", path)
	vObj.SetConfigType("toml")

	err = vObj.ReadRemoteConfig()
	if err != nil {
		err = fmt.Errorf("%v (no section: %v)", err, path)
		return
	}

	x := map[string]interface{}{}

	allKeys := vObj.AllKeys()
	for _, key := range allKeys {
		x[key] = vObj.Get(key)
	}

	result = cast.ToStringMapString(x)

	return
}
