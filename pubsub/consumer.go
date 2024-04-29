package pubsub

import (
	"fmt"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type Message struct {
	Data      []byte
	Key       []byte
	Partition int32
	Offset    int64
}

func ConvertKafkaMessageToMessage(kafkaMsg kafka.Message) Message {
	return Message{
		Data:      kafkaMsg.Value,
		Key:       kafkaMsg.Key,
		Partition: kafkaMsg.TopicPartition.Partition,
		Offset:    int64(kafkaMsg.TopicPartition.Offset),
	}
}

func (m Message) String() string {
	return fmt.Sprintf("Data: %s, Key: %s, Partition: %d, Offset: %d", m.Data, m.Key, m.Partition, m.Offset)
}

type kafkaConsumer struct {
	groupID  string
	consumer *kafka.Consumer
}

func NewKafkaConsumer(brokers string, groupID string) (*kafkaConsumer, error) {
	config := &kafka.ConfigMap{
		"bootstrap.servers": brokers,
		"group.id":          groupID,
		"auto.offset.reset": "earliest",
	}

	c, err := kafka.NewConsumer(config)
	if err != nil {
		return nil, err
	}
	return &kafkaConsumer{groupID: groupID, consumer: c}, nil
}

func (kc *kafkaConsumer) Consume(topic string, errChan chan error) <-chan Message {
	outputChan := make(chan Message)

	topicList := []string{topic}

	err := kc.consumer.SubscribeTopics(topicList, nil)
	if err != nil {
		errChan <- fmt.Errorf("failed to subscribe topics: %v", err)
		close(outputChan)
		return outputChan
	}
	go kc.consumeMessages(errChan, outputChan)

	return outputChan
}

func (kc *kafkaConsumer) consumeMessages(errChan chan error, outputChan chan<- Message) {
	defer close(outputChan)

	for {
		msg, err := kc.consumer.ReadMessage(-1)
		if err != nil {
			errChan <- fmt.Errorf("failed to read message: %v", err)
			continue
		}
		outputChan <- ConvertKafkaMessageToMessage(*msg)
	}
}

func (kc *kafkaConsumer) Close() error {
	return kc.consumer.Close()
}
