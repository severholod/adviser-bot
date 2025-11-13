package main

import (
	"adviser-bot/clients/telegram"
	"adviser-bot/consumer"
	telegramProcessor "adviser-bot/events/telegram"
	"adviser-bot/storage/files"
	"flag"
	"log"
)

const (
	host        = "api.telegram.org"
	storagePath = "files_storage"
	batchSize   = 100
)

// 8356731409:AAFRGQ6fKeihQDsbunuqOTachCKk7-oJVtM
func main() {
	tgClient := telegram.New(host, mustToken())
	eventsProcessor := telegramProcessor.New(tgClient, files.New(storagePath))
	log.Println("Service has been started...")

	eventConsumer := consumer.New(eventsProcessor, eventsProcessor, batchSize)

	if err := eventConsumer.Start(); err != nil {
		log.Fatal("Service is stopped...", err)
	}
}

func mustToken() string {
	token := flag.String("tg-bot-token", "", "access tg-bot token")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is required")
	}

	return *token
}
