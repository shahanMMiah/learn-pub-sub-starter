package main

import (
	"fmt"
	"log"
	"os"

	"github.com/bootdotdev/learn-pub-sub-starter/internal/gamelogic"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/pubsub"
	"github.com/bootdotdev/learn-pub-sub-starter/internal/routing"

	amqp "github.com/rabbitmq/amqp091-go"
)

type serverMsg string

const (
	PAUSEMSG  = "pause"
	RESUMEMSG = "resume"
	QUITMSG   = "quit"
)

func ErrorMsg(pref string, err error) string {
	return fmt.Sprintf("%s %s", pref, err.Error())
}

func main() {
	log.Println("Starting Peril server...")

	connectionString := "amqp://guest:guest@localhost:5672/"

	connection, err := amqp.Dial(connectionString)
	if err != nil {
		log.Println(ErrorMsg("Error starting server..", err))
		os.Exit(1)
	}

	defer connection.Close()
	channel, err := connection.Channel()
	if err != nil {
		log.Println(ErrorMsg("Error starting server..", err))
		os.Exit(1)

	}

	log.Println("Connection to peril server succsefull...")

	gamelogic.PrintServerHelp()

	while := true

	for while {
		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		}

		switch inputs[0] {
		case PAUSEMSG:
			log.Println("Sending Pause Message to channel")
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: true})
			if err != nil {
				log.Println(ErrorMsg("Error sending msg to server..", err))
				os.Exit(1)
			}
		case RESUMEMSG:
			log.Println("Sending Pause Message to channel")
			err = pubsub.PublishJSON(channel, routing.ExchangePerilDirect, routing.PauseKey, routing.PlayingState{IsPaused: false})
			if err != nil {
				log.Println(ErrorMsg("Error sending msg to server..", err))
				os.Exit(1)
			}
		case QUITMSG:
			log.Println("Exiting channel")
			while = false

		default:
			log.Println("Do not know that command....")

		}

	}

	/*
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		<-signalChan
		log.Println("Shutting down server...")
	*/
}
