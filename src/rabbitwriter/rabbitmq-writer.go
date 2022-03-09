package rabbitwriter

import (
	"encoding/json"

	"github.com/MichalMitros/feed-parser/models"
	"github.com/streadway/amqp"
)

type RabbitWriterItnerface interface {
	WriteToQueue(queueName string, shopItems chan models.ShopItem) error
}

type RabbitWriter struct {
	username   string
	password   string
	hostname   string
	connection *amqp.Connection
}

func NewRabbitWriter(
	username string,
	password string,
	hostname string,
) (*RabbitWriter, error) {
	connString := "amqp://" + username + ":" + password + "@" + hostname + "/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		return nil, err
	}

	return &RabbitWriter{
		username:   username,
		password:   password,
		hostname:   hostname,
		connection: conn,
	}, nil
}

func (r RabbitWriter) WriteToQueue(
	queueName string,
	shopItemsInput chan models.ShopItem,
) error {

	ch, err := r.connection.Channel()
	if err != nil {
		return err
	}

	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	go r.rabbitWritingRoutine(
		ch,
		&q,
		shopItemsInput,
	)

	return nil
}

func (r RabbitWriter) rabbitWritingRoutine(
	channel *amqp.Channel,
	queue *amqp.Queue,
	shopItemsInput chan models.ShopItem,
) {
	for item := range shopItemsInput {
		body, _ := json.Marshal(item)
		channel.Publish(
			"",
			queue.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
	}
}
