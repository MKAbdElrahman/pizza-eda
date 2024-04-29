package main

import (
	"encoding/json"
	"fmt"
	"log"
	"pizza/models"
	"pizza/pubsub"
)

func main() {
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
