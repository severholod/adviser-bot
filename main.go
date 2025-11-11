package main

import (
	"flag"
	"log"
)

func main() {
	token := mustToken()
	// tgClient = telegram.New(token)
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
