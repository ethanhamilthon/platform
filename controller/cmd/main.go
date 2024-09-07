package cmd

import (
	"controller/internal/api"
	"controller/internal/config"
	"controller/internal/db"
	"controller/internal/message"
	"controller/internal/services/application"
	"controller/internal/services/auth"
	"controller/internal/services/domain"
	"controller/internal/services/listener"
	"log"
)

func Start() {
	// Create database
	db := db.New(config.DbPath)
	defer db.Close()

	// Create message broker
	broker, err := message.New()
	if err != nil {
		log.Fatal(err)
	}
	// Create stream
  stream,err := broker.CreateStream(message.CreateStramOptions{
		Name: config.ControllerStream,
		Subjects: []string{
			config.ControllerStreamPathPrefix,
		},
	})
	// Create listener
	l := listener.New()
	listener_consumer, err := broker.CreateConsumer(message.CreateConsumerOptions{
		Name:    "controller:listener",
		Subject: config.ControllerAnswer,
	},stream, l.AnswerListener)
	defer listener_consumer.Stop()

	// Create services
	auth := auth.New(db)
	dm := domain.New(db)
	apps := application.New(db, broker, l)

	// Start server
	server := api.New(auth, dm, apps)
	server.Serve()
}
