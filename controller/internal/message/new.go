package message

import (
	"context"
	"controller/internal/config"
	"encoding/json"
	"log"

	"github.com/nats-io/nats.go"
	"github.com/nats-io/nats.go/jetstream"
)

type Broker struct {
	jet jetstream.JetStream
}

func New() (*Broker, error) {
	nc, err := nats.Connect(config.NatsUrl)
	if err != nil {
		return nil, err
	}
	Jet, err := jetstream.New(nc)
	if err != nil {
		return nil, err
	}
	return &Broker{jet: Jet}, nil
}

// CreateWait creates the listener to all messages in the controller stream
// func (b *Broker) createWait(stream string) error {
// 	log.Println("creating waiter")
// 	_, err := b.Consume(stream, stream, func(m jetstream.Msg) {
// 		log.Println("Inside waiter")
// 		b.wait.mu.RLock()
// 		defer b.wait.mu.RUnlock()
// 		if len(b.wait.subs) == 0 {
// 			return
// 		}
// 		log.Println("got data")
// 		id := m.Headers().Get("Nats-Msg-Id")
// 		if ch, ok := b.wait.subs[id]; ok {
// 			ch <- m.Data()
// 			delete(b.wait.subs, id)
// 			m.Ack()
// 		}
// 	}, "waiter1")
// 	return err
// }

// AddWait adds a new listener to controller stream. And returns
// a function to wait for the message. Any message data to the controller stream
// and has a matching Nats-Msg-Id will be returned
// func (b *Broker) AddWait(id string) func() ([]byte, error) {
// 	ch := make(chan []byte, 1)
// 	b.wait.mu.Lock()
// 	b.wait.subs[id] = ch
// 	b.wait.mu.Unlock()
// 	waiterFunc := func() ([]byte, error) {
// 		ticker := time.NewTicker(1 * time.Second)
// 		tries := 3
// 		for {
// 			select {
// 			case msg := <-ch:
// 				return msg, nil
// 			case <-ticker.C:
// 				if tries == 0 {
// 					return nil, errors.New("wait timed out")
// 				}
// 				tries--
// 			}
// 		}
// 	}
// 	return waiterFunc
// }

// CreateStream creates the controller stream
// func (b *Broker) CreateStream(name string, pathprefix string) (jetstream.Stream, error) {
// 	stream, err := b.jet.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
// 		Name:        name,
// 		Description: name + " stream",
// 		Subjects:    []string{pathprefix},
// 	})
// 	return stream, err
// }

// Consume consumes messages from the controller stream
//
//	func (b *Broker) Consume(topic string, stream string, handler func(m jetstream.Msg), durable string) (jetstream.ConsumeContext, error) {
//		con, err := b.jet.CreateOrUpdateConsumer(context.Background(), stream, jetstream.ConsumerConfig{
//			Name:          durable,
//			Durable:       durable,
//			FilterSubject: topic,
//		})
//		if err != nil {
//			return nil, err
//		}
//		consumer, err := con.Consume(func(msg jetstream.Msg) {
//			log.Printf("data for topic: %v, stream:%v, durable:%v, subject:%v", topic, stream, durable, msg.Subject())
//			handler(msg)
//		})
//		return consumer, err
//	}
type CreateStramOptions struct {
	Name     string
	Subjects []string
}

func (b *Broker) CreateStream(options CreateStramOptions) (jetstream.Stream, error) {
	stream, err := b.jet.CreateOrUpdateStream(context.Background(), jetstream.StreamConfig{
		Name:        options.Name,
		Description: options.Name + " stream",
		Subjects:    options.Subjects,
	})
	return stream, err
}

type CreateConsumerOptions struct {
	Name    string
	Subject string
}

func (b *Broker) CreateConsumer(options CreateConsumerOptions, stream jetstream.Stream,
	handler func(m jetstream.Msg)) (jetstream.ConsumeContext, error) {

	consumer, err := stream.CreateOrUpdateConsumer(context.Background(), jetstream.ConsumerConfig{
		Name:          options.Name,
		Durable:       options.Name,
		FilterSubject: options.Subject,
		MaxDeliver:    1,
		AckPolicy:     jetstream.AckExplicitPolicy,
	})
	if err != nil {
		return nil, err
	}
	cs, err := consumer.Consume(handler)
	if err != nil {
		return cs, err
	}
	return cs, err
}

// Publish publishes a message to another straems
func (b *Broker) PublishID(topic string, message interface{}, id string) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}
	log.Printf("publishing to %v with id %v", topic, id)

	_, err = b.jet.PublishMsg(context.Background(), &nats.Msg{
		Data:    body,
		Subject: topic,
	}, jetstream.WithMsgID(id))

	return err
}

func (b *Broker) Publish(topic string, message interface{}) error {
	body, err := json.Marshal(message)
	if err != nil {
		return err
	}

	_, err = b.jet.PublishMsg(context.Background(), &nats.Msg{
		Data:    body,
		Subject: topic,
	})

	return err
}
