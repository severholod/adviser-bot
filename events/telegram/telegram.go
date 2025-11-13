package telegram

import (
	"adviser-bot/clients/telegram"
	"adviser-bot/events"
	"adviser-bot/lib/utils"
	"adviser-bot/storage"
	"errors"
)

type Processor struct {
	tg      *telegram.Client
	offset  int
	storage storage.Storage
}
type Meta struct {
	ChatId   int
	Username string
}

var ErrUnknownEventType = errors.New("Unknown event type")
var ErrUnknownMetaType = errors.New("Unknown meta type")

func New(client *telegram.Client, storage storage.Storage) *Processor {
	return &Processor{
		tg:      client,
		storage: storage,
	}
}

func (p *Processor) Fetch(limit int) ([]events.Event, error) {
	updates, err := p.tg.Updates(p.offset, limit)
	if err != nil {
		return nil, utils.WrapError("can`t get events", err)
	}
	if len(updates) == 0 {
		return nil, nil
	}
	result := make([]events.Event, 0, len(updates))
	for _, update := range updates {
		result = append(result, toEvent(update))
	}

	p.offset = updates[len(updates)-1].ID + 1

	return result, nil
}

func (p *Processor) Process(e events.Event) error {
	switch e.Type {
	case events.Message:
		return p.processMessage(e)
	default:
		return utils.WrapError("can't process message", ErrUnknownEventType)
	}
}

func (p *Processor) processMessage(e events.Event) error {
	meta, err := getEventMeta(e)
	if err != nil {
		return utils.WrapError("can't process message", err)
	}
	if err := p.doCmd(e.Text, meta.ChatId, meta.Username); err != nil {
		return utils.WrapError("can't process message", err)
	}

	return nil
}

func getEventMeta(e events.Event) (Meta, error) {
	res, ok := e.Meta.(Meta)
	if !ok {
		return Meta{}, utils.WrapError("can't process meta", ErrUnknownMetaType)
	}
	return res, nil
}

func toEvent(update telegram.Update) events.Event {
	updateType := fetchType(update)
	event := events.Event{
		Type: updateType,
		Text: fetchText(update),
	}
	if updateType == events.Message {
		event.Meta = Meta{
			ChatId:   update.Message.Chat.Id,
			Username: update.Message.From.Username,
		}
	}
	return event
}
func fetchType(update telegram.Update) events.Type {
	if update.Message == nil {
		return events.Unknown
	}
	return events.Message
}
func fetchText(update telegram.Update) string {
	if update.Message == nil {
		return ""
	}

	return update.Message.Text
}
