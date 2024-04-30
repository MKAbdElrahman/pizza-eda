package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"pizza/models"
	"pizza/pubsub"
	"time"
)

func main() {

	producerErrCh := make(chan error)
	deiveryCh := make(chan pubsub.DeliveryReport)

	producer, err := pubsub.NewConfluentProducer("localhost:9092")

	if err != nil {
		panic(err)
	}

	go func() {
		for e := range producerErrCh {
			fmt.Println(e)
		}
	}()

	go func() {
		for d := range deiveryCh {
			fmt.Println(d)
		}
	}()

	for {

		time.Sleep(time.Second)

		producer.Produce(context.Background(), pubsub.ProducerMessage{
			Topic: "test-topic",
			Key:   []byte("test-key-2"),
			Value: []byte("test-value-2"),
		})

	}

	log.Println("Starting pizza assembler microservice...")

	groupID := "pizza-assemblers-group-1"
	consumer, err := pubsub.NewKafkaConsumer("localhost:9092", groupID)
	if err != nil {
		panic(err)
	}

	errChan := make(chan error)
	go func() {
		for err := range errChan {
			fmt.Println("Error:", err)
		}
	}()

	events := consumer.Consume("pizza-ordered", errChan)

	for e := range events {

		var orderRequest models.PizzaOrder
		if err := json.Unmarshal(e.Data, &orderRequest); err != nil {
			errChan <- err
			continue
		}

		fmt.Println(orderRequest)
	}
}
