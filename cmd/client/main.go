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

const (
	SPWANCMD  = "spawn"
	MOVECMD   = "move"
	STATUSCMD = "status"
	HELPCMD   = "help"
	SPAMCMD   = "spam"
	QUITCMD   = "quit"
)

func inSlice(word string, list []string) bool {
	for _, wrd := range list {
		if wrd == word {
			return true
		}
	}
	return false
}

func ErrorMsg(pref string, err error) string {
	return fmt.Sprintf("%s %s", pref, err.Error())
}

func main() {

	fmt.Println("Starting Peril client...")

	connectionString := "amqp://guest:guest@localhost:5672/"
	connection, err := amqp.Dial(connectionString)
	if err != nil {
		fmt.Println(ErrorMsg("Error starting client..", err))
		os.Exit(1)
	}

	username, err := gamelogic.ClientWelcome()
	if err != nil {
		ErrorMsg("Error client login..", err)
		os.Exit(1)
	}

	pubsub.DeclareAndBind(connection, routing.ExchangePerilDirect, fmt.Sprintf("pause.%s", username), routing.PauseKey, "transient")

	gamestate := gamelogic.NewGameState(username)

	while := true

	for while {

		inputs := gamelogic.GetInput()
		if len(inputs) == 0 {
			continue
		}
		switch inputs[0] {
		case SPWANCMD:
			err = gamestate.CommandSpawn(inputs)
			if err != nil {
				log.Println(ErrorMsg("Error during spawn command..", err))
			}
		case MOVECMD:
			_, err := gamestate.CommandMove(inputs)
			if err != nil {
				log.Println(ErrorMsg("Error during spawn command..", err))
			}
		case STATUSCMD:
			gamestate.CommandStatus()
		case HELPCMD:
			gamelogic.PrintClientHelp()
		case SPAMCMD:
			log.Println("Spamming not allowed yet!")
		case QUITCMD:
			gamelogic.PrintQuit()
			while = false
		default:
			log.Println("error - unkown command..")
		}
	}

	/*
		signalChan := make(chan os.Signal, 1)
		signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM)

		<-signalChan
		fmt.Println("Shutting down client...")
	*/
}
