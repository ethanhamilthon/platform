package listener

import (
	"errors"
	"log"
	"sync"
	"time"

	"github.com/nats-io/nats.go/jetstream"
)

type ListenerService struct {
	mu   *sync.RWMutex
	subs map[string]chan []byte
}

func New() *ListenerService {
	return &ListenerService{
		mu:   &sync.RWMutex{},
		subs: map[string]chan []byte{},
	}
}

func (l *ListenerService) AnswerListener(msg jetstream.Msg) {
	log.Println("message to controller answer listener")
	l.mu.Lock()
	defer l.mu.Unlock()
	if len(l.subs) == 0 {
		return
	}
	id := msg.Headers().Get("Nats-Msg-Id")
	if ch, ok := l.subs[id]; ok {
		ch <- msg.Data()
		delete(l.subs, id)
		msg.Ack()
	}
}

func (l *ListenerService) Add(id string) func() ([]byte, error) {
	ch := make(chan []byte, 1)
	l.mu.Lock()
	l.subs[id] = ch
	l.mu.Unlock()
	waiterFunc := func() ([]byte, error) {
		ticker := time.NewTicker(1 * time.Second)
		tries := 3
		for {
			select {
			case msg := <-ch:
				return msg, nil
			case <-ticker.C:
				if tries == 0 {
					return nil, errors.New("wait timed out")
				}
				tries--
			}
		}
	}
	return waiterFunc
}
