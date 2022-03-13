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
) (*RabbitWriter, error) {
	// Create connection
	connString :=
		"amqp://" +
			options.Username + ":" +
			options.Password + "@" +
			options.Hostname + "/"
	conn, err := amqp.Dial(connString)
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
	errorsOutput chan error,
) error {

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

	// Start new go routine sending items
	// from shopItemsInput channel to the queue
	go r.rabbitWritingRoutine(
		&ch,
		&q,
		shopItemsInput,
		errorsOutput,
	)

	return nil
}

// Writes objects from shopItemsInput
// to the queue in JSON format
func (r *RabbitWriter) rabbitWritingRoutine(
	channel *wabbit.Channel,
	queue *wabbit.Queue,
	shopItemsInput chan models.ShopItem,
	errorsOutput chan error,
) {
	defer close(errorsOutput)
	for item := range shopItemsInput {
		body, _ := json.Marshal(item)
		err := (*channel).Publish(
			"",
			(*queue).Name(),
			body,
			wabbit.Option{},
		)
		// Increment prometheus counter
		if err != nil {
			errorsOutput <- err
			publishedShopItemsFailures.Inc()
		} else {
			publishedShopItems.Inc()
		}
	}
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
