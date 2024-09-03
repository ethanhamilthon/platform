package message

import (
	"time"

	"github.com/nats-io/nats.go"
)

type MessageBroker struct {
	core *nats.Conn
}

func New(nats_addr string) (*MessageBroker, error) {
	mb := &MessageBroker{}
	nc, err := nats.Connect(nats_addr)
	if err != nil {
		return mb, err
	}

	mb.core = nc
	return mb, nil
}

func (m *MessageBroker) Close() {
	m.core.Close()
}

// Response processes messages and returns data to resqued
func (m *MessageBroker) Response(topic string, handler func(data []byte) ([]byte, error)) {
	m.core.Subscribe(topic, func(msg *nats.Msg) {
		res, err := handler(msg.Data)
		if err == nil {
			err = msg.Respond(res)
			if err == nil {
				msg.Ack()
			}
		}
	})
}

// Request sends message to topic, and expects data from consumer
func (m *MessageBroker) Request(topic string, data []byte) ([]byte, error) {
	msg, err := m.core.Request(topic, data, time.Second*2)
	if err != nil {
		return []byte{}, err
	}

	return msg.Data, nil
}
