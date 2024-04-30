package pubsub

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/v2/kafka"
)

type ProducerMessage struct {
	Topic string
	Key   []byte
	Value []byte
}

type DeliveryReport struct {
	Topic     string
	Partition int64
	Offset    int64
}

type Producer interface {
	Produce(c context.Context, msg ProducerMessage) (DeliveryReport, error)
	Close() error
}

type ConfluentProducer struct {
	producer *kafka.Producer
}

func NewConfluentProducer(brokers string) (*ConfluentProducer, error) {
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": brokers})
	if err != nil {
		return nil, err
	}

	cp := &ConfluentProducer{
		producer: p,
	}

	return cp, nil
}

func (p *ConfluentProducer) Close() error {
	p.producer.Close()
	return nil
}

func (p *ConfluentProducer) ProduceAsync(c context.Context, msg ProducerMessage, deliveryReportsCh chan DeliveryReport, errorsCh chan error) {
	go func() {
		deliveryCh := make(chan kafka.Event)
		defer close(deliveryCh)

		err := p.producer.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &msg.Topic, Partition: kafka.PartitionAny},
			Key:            msg.Key,
			Value:          msg.Value,
		}, deliveryCh)
		if err != nil {
			errorsCh <- err
			return
		}

		for e := range deliveryCh {
			select {
			case <-c.Done():
				return
			default:
				switch ev := e.(type) {
				case *kafka.Message:
					if ev.TopicPartition.Error != nil {
						errorsCh <- ev.TopicPartition.Error
					} else {
						deliveryReportsCh <- DeliveryReport{
							Topic:     *ev.TopicPartition.Topic,
							Offset:    int64(ev.TopicPartition.Offset),
							Partition: int64(ev.TopicPartition.Partition),
						}
					}
				}
			}
		}

	}()
}

func (p *ConfluentProducer) Produce(c context.Context, msg ProducerMessage) (DeliveryReport, error) {

	deliveryReportsCh := make(chan DeliveryReport, 1)
	errorsCh := make(chan error, 1)

	p.ProduceAsync(c, msg, deliveryReportsCh, errorsCh)

	select {
	case <-c.Done():
		return DeliveryReport{}, c.Err()
	case err := <-errorsCh:
		return DeliveryReport{}, err
	case deliveryReport := <-deliveryReportsCh:
		return deliveryReport, nil
	}
}
