package rabbitmq

import (
	"log"

	"github.com/streadway/amqp"
)

//HandleConsumedDeliveries handles the consumed deliveries from the queues. Should be called only for a consumer connection
func (c *Connection) HandleConsumedDeliveries(autoAck bool, fn func(Connection, <-chan amqp.Delivery)) {
	delivery, err := c.channel.Consume(c.queue, "", autoAck, false, false, false, nil)
	if err != nil {
		log.Fatalln(err)
	}

	for {
		go fn(*c, delivery)
		if err := <-c.err; err != nil {
			err = c.Reconnect()
			if err != nil {
				log.Println(err)
			}
			delivery, err = c.channel.Consume(c.queue, "", autoAck, false, false, false, nil)
			if err != nil {
				log.Fatalln(err) //raising log fatal if consume fails even after reconnecting
			}
		}
	}
}
