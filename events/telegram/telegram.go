package telegram

import "adviser-bot/clients/telegram"

type Processor struct {
	tg     *telegram.Client
	offset int
	// storage
}

func New(tg *telegram.Client) *Processor {}
