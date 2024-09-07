package builder

import (
	"controller/internal/message"
	"controller/internal/services/listener"
)

type BuilderService struct {
  broker *message.Broker
  listener *listener.ListenerService
}

func New(b *message.Broker, l *listener.ListenerService) *BuilderService {
  return &BuilderService{
    broker: b,
    listener: l,
  }
}


