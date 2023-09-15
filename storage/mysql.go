package storage

import (
	"context"
	"database/sql"
	"log"
	"time"

	"github.com/go-sql-driver/mysql"
	"github.com/win30221/core/config"
)

func GetMysqlDB(path string) (db *sql.DB) {
	var err error

	conf := &mysql.Config{
		Net:                  "tcp",
		Timeout:              time.Second * 3,
		AllowNativePasswords: true,
	}

	defer func() {
		if err != nil {
			log.Fatalf("get mysql error: %s \n - path %s \n - DSN %s", err, path, conf.FormatDSN())
		}
	}()

	conf.Addr, _ = config.GetString(path+"/host", true)
	conf.User, _ = config.GetString(path+"/account", true)
	conf.Passwd, _ = config.GetString(path+"/password", true)
	conf.DBName, _ = config.GetString(path+"/dbname", true)
	conf.Params = map[string]string{"parseTime": "true", "loc": "America/Puerto_Rico"}

	db, err = sql.Open("mysql", conf.FormatDSN())
	if err != nil {
		return
	}

	log.Printf("MySQL connected to `%+v` success", conf.Addr)

	maxOpenConns, _ := config.GetInt(path+"/max_open_conns", true)
	maxIdleConns, _ := config.GetInt(path+"/max_idle_conns", true)
	maxConnLifetime, _ := config.GetDuration(path+"/max_conn_lifetime", true)

	db.SetMaxOpenConns(maxOpenConns)
	db.SetMaxIdleConns(maxIdleConns)
	db.SetConnMaxLifetime(maxConnLifetime)

	conn, err := db.Conn(context.Background())
	if err != nil {
		return
	}
	conn.Close()

	err = db.Ping()
	return
}
