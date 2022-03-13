package rabbitwriter

import (
	"encoding/json"

	"github.com/MichalMitros/feed-parser/models"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
	"github.com/streadway/amqp"
)

// Writer for RabbitMQ queues
type RabbitWriter struct {
	username   string
	password   string
	hostname   string
	connection *amqp.Connection
}

// Connection options for RabibitMQ writer
type RabbitWriterOptions struct {
	Username string
	Password string
	Hostname string
}

// Creates new RabbitWriter instance
func NewRabbitWriter(
	options RabbitWriterOptions,
) (*RabbitWriter, error) {
	writer := RabbitWriter{
		username: options.Username,
		password: options.Password,
		hostname: options.Hostname,
	}

	// Connect
	err := writer.connect()
	if err != nil {
		return nil, err
	}

	return &writer, nil
}

// Creates new connection channel and starts
// new goroutine listening for products in shopItemsInput
// and then sending them to queue queueName
func (r RabbitWriter) WriteToQueue(
	queueName string,
	shopItemsInput chan models.ShopItem,
) error {
	// Declare RabbitMQ queue
	ch, q, err := r.getQueueAndChannel(queueName)
	if err != nil {
		return err
	}

	for item := range shopItemsInput {
		body, _ := json.Marshal(item)
		err = ch.Publish(
			"",
			q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			})
		// Handle connection error
		if err != nil {
			publishedShopItemsFailures.Inc()
			ch, q, err = r.getQueueAndChannel(queueName)
			if err != nil {
				return err
			}
			return err
		} else {
			publishedShopItems.Inc()
		}
	}

	return nil
}

func (r *RabbitWriter) getQueueAndChannel(queueName string) (*amqp.Channel, *amqp.Queue, error) {
	ch, err := r.getChannel()
	if err != nil {
		return nil, nil, err
	}

	// Declare RabbitMQ queue
	q, err := ch.QueueDeclare(
		queueName,
		false,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return nil, nil, err
	}

	return ch, &q, err
}

func (r *RabbitWriter) getChannel() (*amqp.Channel, error) {
	ch, err := r.connection.Channel()
	if err != nil {
		err := r.connect()
		if err != nil {
			return nil, err
		}
		ch, err = r.connection.Channel()
	}

	return ch, err
}

func (r *RabbitWriter) connect() error {
	// Create connection
	connString :=
		"amqp://" +
			r.username + ":" +
			r.password + "@" +
			r.hostname + "/"
	conn, err := amqp.Dial(connString)
	if err != nil {
		return err
	}

	r.connection = conn

	return nil
}

// Prometheus published shop items
var (
	publishedShopItems = promauto.NewCounter(prometheus.CounterOpts{
		Name: "feedparser_rabbitmq_published_items_total",
		Help: "The total number of ShopItems published to RabbitMQ",
	})
	publishedShopItemsFailures = promauto.NewCounter(prometheus.CounterOpts{
		Name: "feedparser_rabbitmq_published_items_failures_total",
		Help: "The total number of failures in publishing ShopItems to RabbitMQ",
	})
)
