package main

import (
	"os"

	log "github.com/sirupsen/logrus"
)

func removeHandler() {
	if len(os.Args) < 3 {
		log.Fatal("Usage: \"remove $subject\"")
	}

	subject := os.Args[2]
	contextLogger := log.WithField("subject", subject)

	err := accountRepo.Remove(os.Args[2])
	if err == nil {
		contextLogger.Info("Removed")
	} else {
		contextLogger.WithError(err).Fatal("Failed to remove")
	}
}
