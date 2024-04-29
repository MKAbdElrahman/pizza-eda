package pubsub

import (
	"errors"
	"fmt"
	"time"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type kafkaProducer struct {
	producer *kafka.Producer
}

func NewKafkaProducer(brokers string) (*kafkaProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil, err
	}
	return &kafkaProducer{producer: p}, nil
}

func (kp *kafkaProducer) Publish(topic string, key, value []byte) error {
	deliveryChan := make(chan kafka.Event)

	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Key:            key,
		Value:          value,
	}

	err := kp.producer.Produce(msg, deliveryChan)
	if err != nil {
		return errors.New("failed to publish message: " + err.Error())
	}

	select {
	case e := <-deliveryChan:
		switch ev := e.(type) {
		case *kafka.Message:
			if ev.TopicPartition.Error != nil {
				return fmt.Errorf("delivery failed: %v", ev.TopicPartition)
			} else {
				fmt.Printf("Delivered message to %v at offset %v\n", ev.TopicPartition, ev.TopicPartition.Offset)
			}
		default:
			fmt.Printf("Unexpected event type received: %T\n", ev)
		}
	case <-time.After(time.Second * 5):
		return errors.New("delivery timed out")
	}

	return nil
}
