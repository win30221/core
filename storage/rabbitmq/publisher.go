package rabbitmq

import (
	"fmt"
	"log"

	"github.com/streadway/amqp"
)

func (c *Connection) Publish(m Message) error {
	select { //non blocking channel - if there is no error will go to default where we do nothing
	case err := <-c.err:
		if err != nil {
			err = c.Reconnect()
			if err != nil {
				log.Println(err)
			}
		}
	default:
	}

	p := amqp.Publishing{
		Headers:       amqp.Table{"type": m.Body.Type},
		ContentType:   m.ContentType,
		CorrelationId: m.CorrelationId,
		Body:          m.Body.Data,
		ReplyTo:       m.ReplyTo,
	}
	if err := c.channel.Publish(c.exchange, m.Queue, false, false, p); err != nil {
		return fmt.Errorf("error in Publishing: %s", err)
	}
	return nil
}
