package storage

import (
	"log"

	"github.com/streadway/amqp"
	"github.com/win30221/core/config"
	"github.com/win30221/core/storage/rabbitmq"
)

type RMQConfig struct {
	Path         string
	Exchange     string
	ExchangeType string
	Queue        string
	Qos          int
}

func GetRabbitMQ(c RMQConfig) (con *rabbitmq.Connection) {
	var queue string

	cfg := amqp.URI{
		Scheme: "amqp",
		Port:   5672,
	}

	cfg.Host, _ = config.GetString(c.Path+"/host", true)
	cfg.Username, _ = config.GetString(c.Path+"/account", true)
	cfg.Password, _ = config.GetString(c.Path+"/password", true)

	if c.Queue != "" {
		queue, _ = config.GetString(c.Path+"/queue."+c.Queue, true)
	}

	con = rabbitmq.NewConnection(cfg, c.Exchange, c.ExchangeType, queue, c.Qos)

	err := con.Reconnect()
	if err != nil {
		log.Fatalf("get rmq error: %s \n - path %s", err, c.Path)
	}

	log.Printf("RMQ connected to `%+v` success", cfg.Host)

	return
}
