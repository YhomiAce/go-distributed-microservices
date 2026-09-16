package events

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"

	amqp "github.com/rabbitmq/amqp091-go"
)

type Consumer struct {
	conn	*amqp.Connection
	queueName  string
}

func NewConsumer(conn *amqp.Connection) (Consumer, error) {
	consumer := Consumer {
		conn: conn,
	}
	err := consumer.setup()
	if err != nil {
		return Consumer{}, err
	}
	return consumer, nil
}

// setup opens a channel and declare the exchange
func (consumer *Consumer) setup() error {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	return declareExchange(channel)
}

// Payload is the type used for pushing event to RabbitMq
type Payload struct {
	Name 	string 	`json:"name"`
	Data 	string 	`json:"data"`
}

// Listen will listen for Queue events
func (consumer *Consumer) Listen(topics []string) error  {
	channel, err := consumer.conn.Channel()
	if err != nil {
		return err
	}
	defer channel.Close()

	q, err := declareRandomQueue(channel)
	if err != nil {
		return err
	}

	for _, topic := range topics {
		err = channel.QueueBind(q.Name, topic, getExchangeName(), false, nil)
		if err != nil {
			log.Println(err.Error())
			return err
		}
	}

	messages, err := channel.Consume(q.Name, "", true, false, false, false, nil)
	if err != nil {
		log.Println(err.Error())
		return err
	}

	forEver := make(chan bool)

	go func() {
		for  message := range messages {
			var payload Payload
			_ = json.Unmarshal(message.Body, &payload)
			handlePayload(payload)
		}
	}()
	log.Printf("[*] Waiting for message [Exchange, Queue][%s, %s]", getExchangeName(), q.Name)
	<- forEver
	return  nil
}

func handlePayload(payload Payload) {
	switch payload.Name {
		case "log", "event":
			err := logEvent(payload)
			if err != nil {
				log.Println(err)
			}
		default:
			err := logEvent(payload)
			if err != nil {
				log.Println(err)
			}
	}
}

func logEvent(entry Payload) error {
	jsonData,_ := json.Marshal(entry)
	var url = "http://logger-service:8082/log"
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		return  err
	}

	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		return  err
	}

	defer res.Body.Close()
	if res.StatusCode != http.StatusAccepted {
		return err
	}
	return nil
}