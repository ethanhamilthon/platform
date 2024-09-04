package cmd

import (
	"balancer/config"
	"balancer/internal/message"
	"balancer/internal/services"
)

func Start() {
	service := services.New()
	msg := message.New(config.NatsUrl)
	msg.AddSubs(config.HttpTopic, service.LaunchHttp)
	msg.AddSubs(config.HttpsTopic, service.LaunchHttps)
	msg.AddSubs(config.AddAppTopic, service.AddApp)
	// Todo: implement add domains
	msg.AddSubs(config.RemoveAppTopic, service.RemoveApp)
	msg.AddSubs(config.ChangeContainerTopic, service.ChangeContainer)
	msg.AddSubs("balancer:ping", func(data []byte) ([]byte, error) {
		return []byte("pong"), nil
	})
}
