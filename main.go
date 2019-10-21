package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

var accountRepo accountRepository

var commandHandlers = map[string]func(){
	"serve":  serveHandler,
	"add":    addHandler,
	"remove": removeHandler,
}

func main() {
	var err error
	accountRepo, err = newFileBasedAccountRepository(accountDirectory)
	if err != nil {
		log.WithField("dir", accountDirectory).WithError(err).Fatal("Failed booting account repository")
	}

	if len(os.Args) < 2 {
		availableCommands := ""
		for command := range commandHandlers {
			availableCommands += " " + command
		}
		log.Fatal("Available commands:", availableCommands)
	}

	command := os.Args[1]
	handler, ok := commandHandlers[command]
	if !ok {
		log.WithField("command", command).Fatal("Unknown command")
	}

	handler()
}
