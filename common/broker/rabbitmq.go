package broker

import (
	"context"
	"fmt"
	"log"
	"time"

	amqp "github.com/rabbitmq/amqp091-go"
	"go.opentelemetry.io/otel"
)

const MaxRetryCount = 3
const DLQ = "dlq_main"

func Connect(user, password, host, port string) (*amqp.Channel, func() error) {
	// Create connection
	address := fmt.Sprintf("amqp://%s:%s@%s:%s", user, password, host, port)

	conn, err := amqp.Dial(address)
	if err != nil {
		log.Fatal(err)
	}

	/// Create channel
	/// Channels are just like multiple connections basically but share a same TCP connection
	ch, err := conn.Channel()
	if err != nil {
		log.Fatal(err)
	}

	/// Orders service -> RabbitMQ -> Payment service
	err = ch.ExchangeDeclare(OrderCreatedEvent, "direct", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	/// Payment Service -> RabbitMQ -> fanout to Orders, Stock, and Kitchen services
	err = ch.ExchangeDeclare(OrderPaidEvent, "fanout", true, false, false, false, nil)
	if err != nil {
		log.Fatal(err)
	}

	/// Dead letter queue
	err = createDLQAndDLX(ch)
	if err != nil {
		log.Fatal(err)
	}

	return ch, conn.Close
}

func HandleRetry(ch *amqp.Channel, d *amqp.Delivery) error {
	if d.Headers == nil {
		d.Headers = amqp.Table{}
	}

	retryCount, ok := d.Headers["x-retry-count"].(int64) /// must be int64
	if !ok {
		retryCount = 0
	}

	retryCount++
	d.Headers["x-retry-count"] = retryCount

	log.Printf("Retrying message %s, retry count: %d", d.Body, retryCount)

	/// React max retry count -> move to dead letters queue
	if retryCount >= MaxRetryCount {
		log.Printf("Moving message to DLQ %s", DLQ)

		/// Publist to DLQ
		return ch.PublishWithContext(context.Background(), "", DLQ, false, false, amqp.Publishing{
			ContentType:  "application/json",
			Headers:      d.Headers,
			Body:         d.Body,
			DeliveryMode: amqp.Persistent,
		})
	}

	/// More number of retries -> more time to wait for retry
	time.Sleep(time.Second * time.Duration(retryCount))

	/// Retry -> publish again
	return ch.PublishWithContext(context.Background(), d.Exchange, d.RoutingKey, false, false, amqp.Publishing{
		ContentType:  "application/json",
		Headers:      d.Headers,
		Body:         d.Body,
		DeliveryMode: amqp.Persistent,
	})
}

// / Create the Dead letter queue
func createDLQAndDLX(ch *amqp.Channel) error {
	q, err := ch.QueueDeclare(
		"main_queue", // name
		true,         // durable
		false,        // delete when unused
		false,        // exclusive
		false,        // no-wait
		nil,          // arguments
	)
	if err != nil {
		return err
	}

	// Declare DLX
	dlx := "dlx_main"
	err = ch.ExchangeDeclare(
		dlx,      // name
		"fanout", // type
		true,     // durable
		false,    // auto-deleted
		false,    // internal
		false,    // no-wait
		nil,      // arguments
	)
	if err != nil {
		return err
	}

	// Bind main queue to DLX
	err = ch.QueueBind(
		q.Name, // queue name
		"",     // routing key
		dlx,    // exchange
		false,
		nil,
	)
	if err != nil {
		return err
	}

	// Declare DLQ
	_, err = ch.QueueDeclare(
		DLQ,   // name
		true,  // durable
		false, // delete when unused
		false, // exclusive
		false, // no-wait
		nil,   // arguments
	)
	if err != nil {
		return err
	}

	return err
}

type AmqpHeaderCarrier map[string]interface{}

func (a AmqpHeaderCarrier) Get(k string) string {
	value, ok := a[k]
	if !ok {
		return ""
	}

	/// type assertion: if value of type string, return value; else panic at runtime
	return value.(string)
}

func (a AmqpHeaderCarrier) Set(k string, v string) {
	a[k] = v
}

func (a AmqpHeaderCarrier) Keys() []string {
	keys := make([]string, len(a))
	i := 0

	for k := range a {
		keys[i] = k
		i++
	}

	return keys
}

func InjectAMQPHeaders(ctx context.Context) map[string]interface{} {
	carrier := make(AmqpHeaderCarrier)
	otel.GetTextMapPropagator().Inject(ctx, carrier) /// inject into context
	return carrier
}

func ExtractAMQPHeader(ctx context.Context, headers map[string]interface{}) context.Context {
	return otel.GetTextMapPropagator().Extract(ctx, AmqpHeaderCarrier(headers)) /// extract concerns from carrier to context
}
