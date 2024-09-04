package config

import "os"

var (
	NatsUrl              string
	HttpTopic            string
	HttpsTopic           string
	AddAppTopic          string
	AddDomainTopic       string
	RemoveAppTopic       string
	ChangeContainerTopic string
)

func loadBrokerConfigs() {
	NatsUrl = getEnv(os.Getenv("NATS_URL"), "nats://localhost:4222")

	HttpTopic = getEnv(os.Getenv("HTTP_TOPIC"), "balancer:launch:http")
	HttpsTopic = getEnv(os.Getenv("HTTPS_TOPIC"), "balancer:launch:https")
	AddAppTopic = getEnv(os.Getenv("ADD_APP_TOPIC"), "balancer:add:app")
	AddDomainTopic = getEnv(os.Getenv("ADD_DOMAIN_TOPIC"), "balancer:add:domain")
	RemoveAppTopic = getEnv(os.Getenv("REMOVE_APP_TOPIC"), "balancer:remove:app")
	ChangeContainerTopic = getEnv(os.Getenv("CHANGE_CONTAINER_TOPIC"), "balancer:change:container")
}
