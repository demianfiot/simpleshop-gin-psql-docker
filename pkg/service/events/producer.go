package events

import (
	"context"
	"encoding/json"
	"prac/todo"
	"strconv"

	kafka "github.com/segmentio/kafka-go"
)

type Producer interface {
	PublishOrderCreated(ctx context.Context, event todo.OrderCreatedEvent) error
} //podia na zamovlem9

//	type Producer interface {
//		Publish(ctx context.Context, event Event) error
//	} // OrderCreatedEvent in Event
type KafkaProducer struct {
	writer *kafka.Writer // writer pidluchenn9 + vidprav povidom
}

func NewKafkaProducer(brokers []string) *KafkaProducer {
	return &KafkaProducer{
		writer: &kafka.Writer{
			Addr:     kafka.TCP(brokers...), // pidluchenn9 do brokeriv
			Topic:    "order_events",        // tema dlya vidpravki povidomlen
			Balancer: &kafka.LeastBytes{},   // balancer dlya rozpodilu povidomlen
		},
	} // stvorenn9 novogo KafkaProducer
}
func (p *KafkaProducer) PublishOrderCreated(
	ctx context.Context,
	event todo.OrderCreatedEvent,
) error {

	payload, err := json.Marshal(event)
	if err != nil {
		return err
	}

	return p.writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(strconv.FormatInt(int64(event.UserID), 10)), // -> por9dok diy dl9 1 korustyvacha
		Value: payload,
	})
}
