package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func addHandler() {
	if len(os.Args) < 4 {
		log.Fatal("Usage: \"add $subject $password\"")
	}

	subject := os.Args[2]
	password := os.Args[3]
	contextLogger := log.WithField("subject", subject)

	err := accountRepo.Add(subject, password)
	if err == nil {
		contextLogger.Info("Added")
	} else {
		contextLogger.WithError(err).Fatal("Failed to add")
	}
}
