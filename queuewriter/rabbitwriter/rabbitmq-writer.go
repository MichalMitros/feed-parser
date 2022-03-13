package rabbitwriter

import (
	"encoding/json"

	"github.com/MichalMitros/feed-parser/models"
	"github.com/NeowayLabs/wabbit"
	"github.com/NeowayLabs/wabbit/amqp"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Writer for RabbitMQ queues
type RabbitWriter struct {
	username   string
	password   string
	hostname   string
	connection *amqp.Conn
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
	dialFunc func(uri string) (*amqp.Conn, error),
) (*RabbitWriter, error) {
	// Create connection
	connString :=
		"amqp://" +
			options.Username + ":" +
			options.Password + "@" +
			options.Hostname + "/"
	conn, err := dialFunc(connString)
	if err != nil {
		return nil, err
	}

	return &RabbitWriter{
		username:   options.Username,
		password:   options.Password,
		hostname:   options.Hostname,
		connection: conn,
	}, nil
}

// Creates new connection channel and starts
// new goroutine listening for products in shopItemsInput
// and then sending them to queue queueName
func (r RabbitWriter) WriteToQueue(
	queueName string,
	shopItemsInput chan models.ShopItem,
) error {
	var err error
	// Get channel from connection
	ch, err := r.connection.Channel()
	if err != nil {
		return err
	}

	// Declare RabbitMQ queue
	q, err := ch.QueueDeclare(
		queueName,
		wabbit.Option{},
	)
	if err != nil {
		return err
	}

	for item := range shopItemsInput {
		body, _ := json.Marshal(item)
		err := (ch).Publish(
			"",
			q.Name(),
			body,
			wabbit.Option{},
		)
		// Increment prometheus counter
		if err != nil {
			publishedShopItemsFailures.Inc()
			return err
		} else {
			publishedShopItems.Inc()
		}
	}

	return err
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
