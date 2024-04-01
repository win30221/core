package storage

import (
	"context"
	"log"
	"net/url"
	"strings"
	"time"

	"github.com/win30221/core/config"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func GetMongoDB(path string, rp *readpref.ReadPref) (db *mongo.Client) {
	var err error
	defer func() {
		if err != nil {
			log.Fatalf("get mongo error: %s \n - path %s", err, path)
		}
	}()

	dbHost, _ := config.GetStringSlice(path+"/host", true)
	hosts := strings.Join(dbHost, ",")

	dbAccount, _ := config.GetString(path+"/account", true)
	dbPassword, _ := config.GetString(path+"/password", true)
	dbPoolSize, _ := config.GetInt(path+"/pool_size", true)

	encodedAccount := url.QueryEscape(dbAccount)
	encodedPassword := url.QueryEscape(dbPassword)

	clientOptions := options.Client().ApplyURI("mongodb://" + encodedAccount + ":" + encodedPassword + "@" + hosts).SetMaxPoolSize(uint64(dbPoolSize))
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	if rp != nil {
		clientOptions.SetReadPreference(rp)
	}

	// Connect to MongoDB
	db, err = mongo.Connect(ctx, clientOptions)
	if err != nil {
		return
	}

	log.Printf("Mongo connected to `%+v` success", hosts)

	// Check the connection
	err = db.Ping(ctx, readpref.Primary())

	return
}
