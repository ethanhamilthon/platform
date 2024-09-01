package main

import (
	"balancer/internal/message"
	"balancer/internal/services"
)

func main() {
	service := services.New()
	msg := message.New("nats://localhost:4222/")

	msg.AddSubs("balancer:launch:http", service.LaunchHttp)
	msg.AddSubs("balancer:launch:https", service.LaunchHttps)
	msg.AddSubs("balancer:add:app", service.AddApp)
	// Todo: implement add domains
	msg.AddSubs("balancer:remove:app", service.RemoveApp)
	msg.AddSubs("balancer:change:container", service.ChangeContainer)

	select {}
}
