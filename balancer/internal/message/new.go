package message

import (
	"balancer/internal/utils"
	"log"

	"github.com/nats-io/nats.go"
)

type Message struct {
	client *nats.Conn
}

func New(addr string) *Message {
	nc, err := nats.Connect(addr)
	if err != nil {
		log.Fatalln("failed to connect NATS")
	}
	return &Message{
		client: nc,
	}
}

// Adds new subscriber to NATS server
func (m *Message) AddSubs(topic string, callback func(data []byte) ([]byte, error)) {
	m.client.Subscribe(topic, func(msg *nats.Msg) {
		answer, err := callback(msg.Data)
		if err != nil {
			msg.Respond(utils.Error(err))
			return
		}
		msg.Respond(answer)
	})
}
