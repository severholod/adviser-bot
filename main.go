package main

import (
	"adviser-bot/clients/telegram"
	"flag"
	"log"
)

const host = "api.telegram.org"

func main() {
	tgClient := telegram.New(host, mustToken())
	// fetcher = fetcher.New()
	// processor = processor.New()

	// consumer.Start(fetcher, processor)
}

func mustToken() string {
	token := flag.String("token-bot-token", "", "access tg-bot token")
	flag.Parse()

	if *token == "" {
		log.Fatal("token is required")
	}

	return *token
}
